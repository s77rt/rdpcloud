<?php

if (!defined("WHMCS")) {
	die("This file cannot be accessed directly");
}

require dirname(__FILE__) . '/vendor/autoload.php';

function serverConfigCheck(array $params) {
	$server = $params["server"];
	$serversecure = $params["serversecure"];

	if ($server !== true) {
		return "No server is assigned";
	}

	if ($serversecure !== true) {
		return "SSL/TLS connection must be enabled in the Server Configuration";
	}

	return null;
}

function getOpts() {
	return [
		"credentials" => Grpc\ChannelCredentials::createInsecure(),
	];
}

function getToken(string $addr, string $username, string $password) {
	$client = new Services\Secauthn\SecauthnClient($addr, getOpts());

	$user = new Models\Secauthn\User_3();
	$user->setUsername($username);
	$user->setPassword($password);

	$request = new Services\Secauthn\LogonUserRequest();
	$request->setUser($user);

	list($response, $status) = $client->LogonUser($request)->wait();

	if ($status->code !== Grpc\STATUS_OK) {
		return [null, "ERROR [LogonUser]: " . $status->code . ", " . $status->details];
	}

	return [$response->getToken(), null];
}

// $size must be UPPERCASE
function parseSize(string $size) {
	static $units = array(
		"B" => 1,
		"KB" => 1024,
		"MB" => 1048576,
		"GB" => 1073741824,
		"TB" => 1099511627776,
	);

	if ($size === "UNLIMITED") {
		return -1;
	}

	$number = floatval($size);
	$unit = preg_replace("/[^A-Z]+/", "", $size);

	return intval($number * $units[$unit]);
}

function parseConfigOptions(array $params) {
	$localGroups = array_filter(array_map('trim', explode(",", $params["configoption1"])));
	$quotaVolumes = array_filter(array_map('trim', explode(",", $params["configoption2"])));
	$quotaVolumesThresholds = array_filter(array_map("trim", explode(",", $params["configoption3"])));
	$quotaVolumesLimits = array_filter(array_map('trim', explode(",", $params["configoption4"])));

	if (count($quotaVolumes) !== count($quotaVolumesThresholds)) {
		return [null, null, null, null, "Number of entries in Quota Volumes and Quota Volumes Thresholds mismatch"];
	}

	if (count($quotaVolumes) !== count($quotaVolumesLimits)) {
		return [null, null, null, null, "Number of entries in Quota Volumes and Quota Volumes Limits mismatch"];
	}

	$quotaVolumesThresholds = array_map("strtoupper", $quotaVolumesThresholds);
	$quotaVolumesThresholds = array_map("parseSize", $quotaVolumesThresholds);
	$quotaVolumesLimits = array_map("strtoupper", $quotaVolumesLimits);
	$quotaVolumesLimits = array_map("parseSize", $quotaVolumesLimits);

	if (in_array(0, $quotaVolumesThresholds)) {
		return [null, null, null, null, "Quota Volumes Thresholds contains an invalid value"];
	}

	if (in_array(0, $quotaVolumesLimits)) {
		return [null, null, null, null, "Quota Volumes Limits contains an invalid value"];
	}

	return [$localGroups, $quotaVolumes, $quotaVolumesThresholds, $quotaVolumesLimits, null];
}

// genUsernameAndPassword generates a username and a password that meets the Windows password policy requirements
// a new username will be generared only if the current username is not set (empty) or if $force_gen_username is set to true
function genUsernameAndPassword(array $params, bool $force_gen_username = false) {
	$username = $params["username"];
	if (strlen($username) === 0 || $force_gen_username === true) {
		$username = "tbd";
	}

	$KEYSPACE_UPPER_CASE = "ABCDEFGHIJKLMNOPQRSTUVWXYZ";
	$KEYSPACE_UPPER_CASE__LENGTH = strlen($KEYSPACE_UPPER_CASE);
	$KEYSPACE_LOWER_CASE = "abcdefghijklmnopqrstuvwxyz";
	$KEYSPACE_LOWER_CASE__LENGTH = strlen($KEYSPACE_LOWER_CASE);
	$KEYSPACE_NUMBERS = "0123456789";
	$KEYSPACE_NUMBERS__LENGTH = strlen($KEYSPACE_NUMBERS);
	$KEYSPACE_SYMBOLS = "~!@#$%^&*_-+=`|\\(){}[]:;\"'<>,.?/";
	$KEYSPACE_SYMBOLS = "~!@#$%^&*_-+=`|(){}[]:;'<>,.?/"; // TMP
	$KEYSPACE_SYMBOLS__LENGTH = strlen($KEYSPACE_SYMBOLS);

	do {
		$keyspace = "";
		for ($i = 0; $i < random_int(2, 8); ++$i) {
			$keyspace .= $KEYSPACE_UPPER_CASE[random_int(0, $KEYSPACE_UPPER_CASE__LENGTH - 1)];
		}
		for ($i = 0; $i < random_int(2, 8); ++$i) {
			$keyspace .= $KEYSPACE_LOWER_CASE[random_int(0, $KEYSPACE_LOWER_CASE__LENGTH - 1)];
		}
		for ($i = 0; $i < random_int(2, 8); ++$i) {
			$keyspace .= $KEYSPACE_NUMBERS[random_int(0, $KEYSPACE_NUMBERS__LENGTH - 1)];
		}
		for ($i = 0; $i < random_int(2, 8); ++$i) {
			$keyspace .= $KEYSPACE_SYMBOLS[random_int(0, $KEYSPACE_SYMBOLS__LENGTH - 1)];
		}

		$positions = range(0, strlen($keyspace) - 1);

		$password = "";
		for ($i = 0; $i < strlen($keyspace); ++$i) {
			$random_position = key(array_slice($positions, random_int(0, count($positions) - 1), 1, true));
			unset($positions[$random_position]);
			$password[$i] = $keyspace[$random_position];
		}
	} while (strlen($username) >= 3 && stripos($password, $username) !== false);

	return [$username, $password];
}

function rdpcloud_MetaData() {
	return array(
		'DisplayName' => 'RDPCloud',
		'APIVersion' => '1.1',
		'RequiresServer' => true,
		'DefaultNonSSLPort' => '5027',
		'DefaultSSLPort' => '5027',
		'ServiceSingleSignOnLabel' => 'Access Panel',
		'AdminSingleSignOnLabel' => 'Access Panel',
	);
}

function rdpcloud_ConfigOptions() {
	return array(
		'Local Groups' => array(
			'Type' => 'text',
			'Description' => 'Local groups (comma separated)',
			'Default' => 'Users, Remote Desktop Users',
		),
		'Quota Volumes' => array(
			'Type' => 'text',
			'Description' => 'Quota volumes (comma separated). Example: C:\\, D:\\. Leave empty if you don\'t want to set quota',
		),
		'Quota Volumes Thresholds' => array(
			'Type' => 'text',
			'Description' => 'Quota volumes thresholds (comma separated). Example 1: 150MB, 1.75GB. Example 2: 150MB, UNLIMITED',
		),
		'Quota Volumes Limits' => array(
			'Type' => 'text',
			'Description' => 'Quota volumes limits (comma separated). Example 1: 200MB, 2GB. Example 2: 150MB, UNLIMITED',
		),
	);
}

function rdpcloud_CreateAccount(array $params) {
	try {
		$error = serverConfigCheck($params);
		if (!is_null($error)) {
			throw new Exception($error);
		}

		$addr = $params["serverip"] . ":" . $params["serverport"];

		list($token, $error) = getToken($addr, $params["serverusername"], $params["serverpassword"]);
		if (!is_null($error)) {
			throw new Exception($error);
		}

		list($localGroups, $quotaVolumes, $quotaVolumesThresholds, $quotaVolumesLimits, $error) = parseConfigOptions($params);
		if (!is_null($error)) {
			throw new Exception($error);
		}

		list($username, $password) = genUsernameAndPassword($params, false);
		$params["model"]->serviceProperties->save(["Username" => $username, "Password" => $password]);

		$netmgmtClient = new Services\Netmgmt\NetmgmtClient($addr, getOpts());

		$user = new Models\Netmgmt\User_3();
		$user->setUsername($username);
		$user->setPassword($password);

		$request = new \Services\Netmgmt\AddUserRequest();
		$request->setUser($user);

		list($response, $status) = $netmgmtClient->AddUser($request, ["authorization" => ["Bearer " . $token]])->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			throw new Exception("ERROR [AddUser]: " . $status->code . ", " . $status->details);
		}

		foreach ($localGroups as $localGroup_Groupname) {
			$user = new Models\Netmgmt\User_1();
			$user->setUsername($username);

			$localGroup = new Models\Netmgmt\LocalGroup_1();
			$localGroup->setGroupname($localGroup_Groupname);

			$request = new \Services\Netmgmt\AddUserToLocalGroupRequest();
			$request->setUser($user);
			$request->setLocalGroup($localGroup);

			list($response, $status) = $netmgmtClient->AddUserToLocalGroup($request, ["authorization" => ["Bearer " . $token]])->wait();
			if ($status->code !== Grpc\STATUS_OK) {
				throw new Exception("User created but an error occurred. ERROR [AddUserToLocalGroup]: " . $status->code . ", " . $status->details . ". IMPORTANT: DELETE THE USER BEFORE TRYING AGAIN");
			}
		}

		$fileioClient = new Services\Fileio\FileioClient($addr, getOpts());
		foreach ($quotaVolumes as $key => $quotaVolume_Volumepath) {
			$volumePath = $quotaVolume_Volumepath;

			$user = new Models\Fileio\User_1();
			$user->setUsername($username);

			$quotaEntry = new Models\Fileio\QuotaEntry_6();
			$quotaEntry->setQuotaThreshold($quotaVolumesThresholds[$key]);
			$quotaEntry->setQuotaLimit($quotaVolumesLimits[$key]);

			$request = new \Services\Fileio\SetUserQuotaEntryRequest();
			$request->setVolumePath($volumePath);
			$request->setUser($user);
			$request->setQuotaEntry($quotaEntry);

			list($response, $status) = $fileioClient->SetUserQuotaEntry($request, ["authorization" => ["Bearer " . $token]])->wait();
			if ($status->code !== Grpc\STATUS_OK) {
				throw new Exception("User created but an error occurred. ERROR [SetUserQuotaEntry]: " . $status->code . ", " . $status->details . ". IMPORTANT: DELETE THE USER BEFORE TRYING AGAIN");
			}
		}

	} catch (Exception $e) {
		// Record the error in WHMCS's module log.
		logModuleCall(
			"rdpcloud",
			__FUNCTION__,
			$params,
			$e->getMessage(),
			$e->getTraceAsString()
		);

		return $e->getMessage();
	}

	return "success";
}

function rdpcloud_TerminateAccount(array $params) {
	return "success";
}
