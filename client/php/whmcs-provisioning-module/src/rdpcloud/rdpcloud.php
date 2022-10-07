<?php

if (!defined("WHMCS")) {
	die("This file cannot be accessed directly");
}

require dirname(__FILE__) . '/vendor/autoload.php';

function getOpts() {
	return [
		"credentials" => Grpc\ChannelCredentials::createInsecure(),
	];
}

function parseSize(string $size) {
	static $units = array(
		"B" => 1,
		"KB" => 1024, "KIB" => 1024,
		"MB" => 1048576, "MIB" => 1048576,
		"GB" => 1073741824, "GIB" => 1073741824,
		"TB" => 1099511627776, "TIB" => 1099511627776,
	);

	$size = strtoupper($size);

	if ($size === "UNLIMITED") {
		return -1;
	}

	$number = floatval($size);
	$unit = preg_replace("/[^A-Z]+/", "", $size);

	return intval($number * $units[$unit]);
}

function parseConfigOptions(array $params) {
	$localGroupsNames = array_filter(array_map("trim", explode(",", $params["configoption1"])));
	$quotaVolumes = array_filter(array_map("trim", explode(",", $params["configoption2"])));
	$quotaVolumesThresholds = array_filter(array_map("trim", explode(",", $params["configoption3"])));
	$quotaVolumesLimits = array_filter(array_map("trim", explode(",", $params["configoption4"])));

	if (count($quotaVolumes) !== count($quotaVolumesThresholds)) {
		return [null, null, null, null, "Number of entries in Quota Volumes and Quota Volumes Thresholds mismatch"];
	}

	if (count($quotaVolumes) !== count($quotaVolumesLimits)) {
		return [null, null, null, null, "Number of entries in Quota Volumes and Quota Volumes Limits mismatch"];
	}

	$quotaVolumesThresholds = array_map("parseSize", $quotaVolumesThresholds);
	$quotaVolumesLimits = array_map("parseSize", $quotaVolumesLimits);

	if (in_array(0, $quotaVolumesThresholds)) {
		return [null, null, null, null, "Quota Volumes Thresholds contains an invalid value"];
	}

	if (in_array(0, $quotaVolumesLimits)) {
		return [null, null, null, null, "Quota Volumes Limits contains an invalid value"];
	}

	return [$localGroupsNames, $quotaVolumes, $quotaVolumesThresholds, $quotaVolumesLimits, null];
}

// genUsernameAndPassword generates a username and a password that meets the Windows password policy requirements
// a new username will be generared only if the current username is not set (empty) or if $username_level is greater than zero
function genUsernameAndPassword(array $params, int $username_level = 0) {
	$username = $params["username"];
	$password = $params["password"];

	if (strlen($username) === 0 || $username_level > 0) {
		if ($username_level == 0) {
			$username_level = 1;
		}

		$firstname = ucfirst(preg_replace("/[^a-zA-Z0-9]+/", "", $params["clientsdetails"]["firstname"]));
		$lastname = ucfirst(preg_replace("/[^a-zA-Z0-9]+/", "", $params["clientsdetails"]["lastname"]));
		$emailname = "";
		if (filter_var($params["clientsdetails"]["email"], FILTER_VALIDATE_EMAIL) !== false) {
			$email = $params["clientsdetails"]["email"];
			$emailname = substr($email, 0, strrpos($email, "@"));
			$emailname = ucfirst(preg_replace("/[^a-zA-Z0-9]+/", "", $emailname));
		}

		$USERNAME_MIN_LENGTH = 7;
		$USERNAME_MAX_LENGTH = 17;

		$username = "";

		if ($username_level <= 3) {
			if (strlen($username) < $USERNAME_MIN_LENGTH || $username_level >= 1) {
				$username = $firstname . $username;
			}
			if (strlen($username) < $USERNAME_MIN_LENGTH || $username_level >= 2) {
				$username = $lastname . $username;
			}
			if (strlen($username) < $USERNAME_MIN_LENGTH || $username_level >= 3) {
				$username = $emailname . $username;
			}

			if (strlen($username) < $USERNAME_MIN_LENGTH) {
				$missing_length = $USERNAME_MIN_LENGTH - strlen($username);
				$username .= random_int(pow(10, $missing_length - 1), pow(10, $missing_length) - 1);
			}

			$username = substr($username, 0, $USERNAME_MAX_LENGTH - 2) . random_int(10, 99);
		} else {
			$username = "User" . random_int(1000, 9999);
		}

		$username = substr($username, 0, $USERNAME_MAX_LENGTH);
	}

	$old_password = $params["password"];

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
	} while ((strlen($username) >= 3 && stripos($password, $username) !== false) || (strlen($old_password) >= 3 && stripos($password, $old_password) !== false));

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
		// Check server config
		if ($params["server"] !== true) {
			throw new Exception("No server is assigned");
		}

		if ($params["serversecure"] !== true) {
			throw new Exception("SSL/TLS connection must be enabled in the Server Configuration");
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Parse parameters
		$addr = $params["serverip"] . ":" . $params["serverport"];

		$serverusername = $params["serverusername"];
		$serverpassword = $params["serverpassword"];

		$username = $params["username"];
		$password = $params["password"];

		list($localGroupsNames, $quotaVolumes, $quotaVolumesThresholds, $quotaVolumesLimits, $error) = parseConfigOptions($params);
		if (!is_null($error)) {
			throw new Exception($error);
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Create clients
		$opts = getOpts();

		$secauthnClient = new Services\Secauthn\SecauthnClient($addr, $opts);
		$netmgmtClient = new Services\Netmgmt\NetmgmtClient($addr, $opts);
		$fileioClient = new Services\Fileio\FileioClient($addr, $opts);
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Get token
		$user = new Models\Secauthn\User_3();
		$user->setUsername($serverusername);
		$user->setPassword($serverpassword);

		$request = new Services\Secauthn\LogonUserRequest();
		$request->setUser($user);

		list($response, $status) = $secauthnClient->LogonUser($request)->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			throw new Exception("Unable to get token: " . $status->code . ", " . $status->details);
		}

		$token = $response->getToken();
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Add user
		list($username, $password) = genUsernameAndPassword($params, 0);
		$params["model"]->serviceProperties->save(["Username" => $username, "Password" => $password]);

		$user = new Models\Netmgmt\User_3();
		$user->setUsername($username);
		$user->setPassword($password);

		$request = new Services\Netmgmt\AddUserRequest();
		$request->setUser($user);

		list($response, $status) = $netmgmtClient->AddUser($request, ["authorization" => ["Bearer " . $token]])->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			throw new Exception("ERROR [AddUser]: ({$status->code}) {$status->details}");
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Add user to local groups
		foreach ($localGroupsNames as $localGroupName) {
			$groupname = $localGroupName;

			$user = new Models\Netmgmt\User_1();
			$user->setUsername($username);

			$localGroup = new Models\Netmgmt\LocalGroup_1();
			$localGroup->setGroupname($groupname);

			$request = new Services\Netmgmt\AddUserToLocalGroupRequest();
			$request->setUser($user);
			$request->setLocalGroup($localGroup);

			list($response, $status) = $netmgmtClient->AddUserToLocalGroup($request, ["authorization" => ["Bearer " . $token]])->wait();
			if ($status->code !== Grpc\STATUS_OK) {
				throw new Exception("User created but an error occurred. ERROR [AddUserToLocalGroup] [{$groupname}]: ({$status->code}) {$status->details}. IMPORTANT: DELETE THE USER BEFORE TRYING AGAIN");
			}
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Set user quota entries
		foreach ($quotaVolumes as $key => $volumePath) {
			$user = new Models\Fileio\User_1();
			$user->setUsername($username);

			$quotaEntry = new Models\Fileio\QuotaEntry_6();
			$quotaEntry->setQuotaThreshold($quotaVolumesThresholds[$key]);
			$quotaEntry->setQuotaLimit($quotaVolumesLimits[$key]);

			$request = new Services\Fileio\SetUserQuotaEntryRequest();
			$request->setVolumePath($volumePath);
			$request->setUser($user);
			$request->setQuotaEntry($quotaEntry);

			list($response, $status) = $fileioClient->SetUserQuotaEntry($request, ["authorization" => ["Bearer " . $token]])->wait();
			if ($status->code !== Grpc\STATUS_OK) {
				throw new Exception("User created but an error occurred. ERROR [SetUserQuotaEntry] [{$volumePath}]: ({$status->code}) {$status->details}. IMPORTANT: DELETE THE USER BEFORE TRYING AGAIN");
			}
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	} catch (Exception $e) {
		// Record the error in WHMCS's module log.
		logModuleCall(
			"rdpcloud",
			__FUNCTION__,
			$params,
			$e->getMessage(),
			$e->getTraceAsString()
		);
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		return $e->getMessage();
	}

	return "success";
}

function rdpcloud_TerminateAccount(array $params) {
	try {
		// Check server config
		if ($params["server"] !== true) {
			throw new Exception("No server is assigned");
		}

		if ($params["serversecure"] !== true) {
			throw new Exception("SSL/TLS connection must be enabled in the Server Configuration");
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Parse parameters
		$addr = $params["serverip"] . ":" . $params["serverport"];

		$serverusername = $params["serverusername"];
		$serverpassword = $params["serverpassword"];

		$username = $params["username"];
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Create clients
		$opts = getOpts();

		$secauthnClient = new Services\Secauthn\SecauthnClient($addr, $opts);
		$netmgmtClient = new Services\Netmgmt\NetmgmtClient($addr, $opts);
		$termservClient = new Services\Termserv\TermservClient($addr, $opts);
		$secauthzClient = new Services\Secauthz\SecauthzClient($addr, $opts);
		$shellClient = new Services\Shell\ShellClient($addr, $opts);
		$fileioClient = new Services\Fileio\FileioClient($addr, $opts);
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Get token
		$user = new Models\Secauthn\User_3();
		$user->setUsername($serverusername);
		$user->setPassword($serverpassword);

		$request = new Services\Secauthn\LogonUserRequest();
		$request->setUser($user);

		list($response, $status) = $secauthnClient->LogonUser($request)->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			throw new Exception("Unable to get token: " . $status->code . ", " . $status->details);
		}

		$token = $response->getToken();
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Disable user
		$user = new Models\Netmgmt\User_1();
		$user->setUsername($username);

		$request = new Services\Netmgmt\DisableUserRequest();
		$request->setUser($user);

		list($response, $status) = $netmgmtClient->DisableUser($request, ["authorization" => ["Bearer " . $token]])->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			if ($status->code === Grpc\STATUS_FAILED_PRECONDITION && $status->details === "User is already disabled") {
				;
			} else {
				throw new Exception("ERROR [DisableUser]: " . $status->code . ", " . $status->details);
			}
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Logoff user
		$user = new Models\Termserv\User_1();
		$user->setUsername($username);

		$request = new Services\Termserv\LogoffUserRequest();
		$request->setUser($user);

		list($response, $status) = $termservClient->LogoffUser($request, ["authorization" => ["Bearer " . $token]])->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			if ($status->code === Grpc\STATUS_NOT_FOUND && $status->details === "User not found / User not logged in") {
				;
			} else {
				throw new Exception("ERROR [LogoffUser]: " . $status->code . ", " . $status->details);
			}
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Get user sid
		$request = new Services\Secauthz\LookupAccountSidByUsernameRequest();
		$request->setUsername($username);

		list($response, $status) = $secauthzClient->LookupAccountSidByUsername($request, ["authorization" => ["Bearer " . $token]])->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			throw new Exception("ERROR [LookupAccountSidByUsername]: " . $status->code . ", " . $status->details);
		}

		$sid = $response->getSid();
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Delete user profile
		$request = new Services\Shell\DeleteProfileRequest();
		$request->setSid($sid);

		list($response, $status) = $shellClient->DeleteProfile($request, ["authorization" => ["Bearer " . $token]])->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			// This may fail if the user don't have a profile yet (didn't login yet)
			; // Silent failure. Move on...
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Get volumes paths (one volume path at most for each volume)
		$request = new Services\Fileio\GetVolumesRequest();

		list($response, $status) = $fileioClient->GetVolumes($request, ["authorization" => ["Bearer " . $token]])->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			throw new Exception("ERROR [GetVolumes]: " . $status->code . ", " . $status->details);
		}

		$volumesPaths = [];

		$volumes = $response->getVolumes();
		foreach ($volumes as $volume) {
			$volumePaths = $volume->getPaths();
			if (count($volumePaths) > 0) {
				array_push($volumesPaths, $volumePaths[0]);
			}
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Delete user quota entries
		foreach ($volumesPaths as $volumePath) {
			$user = new Models\Fileio\User_1();
			$user->setUsername($username);

			$request = new Services\Fileio\DeleteUserQuotaEntryRequest();
			$request->setVolumePath($volumePath);
			$request->setUser($user);

			list($response, $status) = $fileioClient->DeleteUserQuotaEntry($request, ["authorization" => ["Bearer " . $token]])->wait();
			if ($status->code !== Grpc\STATUS_OK) {
				throw new Exception("ERROR [DeleteUserQuotaEntry] [{$volumePath}]: ({$status->code}) {$status->details}");
			}
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Delete user
		$user = new Models\Netmgmt\User_1();
		$user->setUsername($username);

		$request = new Services\Netmgmt\DeleteUserRequest();
		$request->setUser($user);

		list($response, $status) = $netmgmtClient->DeleteUser($request, ["authorization" => ["Bearer " . $token]])->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			throw new Exception("ERROR [DeleteUser]: " . $status->code . ", " . $status->details);
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	} catch (Exception $e) {
		// Record the error in WHMCS's module log.
		logModuleCall(
			"rdpcloud",
			__FUNCTION__,
			$params,
			$e->getMessage(),
			$e->getTraceAsString()
		);
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		return $e->getMessage();
	}

	return "success";
}

function rdpcloud_ChangePackage(array $params) {
	try {
		// Check server config
		if ($params["server"] !== true) {
			throw new Exception("No server is assigned");
		}

		if ($params["serversecure"] !== true) {
			throw new Exception("SSL/TLS connection must be enabled in the Server Configuration");
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Parse parameters
		$addr = $params["serverip"] . ":" . $params["serverport"];

		$serverusername = $params["serverusername"];
		$serverpassword = $params["serverpassword"];

		$username = $params["username"];

		list($localGroupsNames, $quotaVolumes, $quotaVolumesThresholds, $quotaVolumesLimits, $error) = parseConfigOptions($params);
		if (!is_null($error)) {
			throw new Exception($error);
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Create clients
		$opts = getOpts();

		$secauthnClient = new Services\Secauthn\SecauthnClient($addr, $opts);
		$netmgmtClient = new Services\Netmgmt\NetmgmtClient($addr, $opts);
		$fileioClient = new Services\Fileio\FileioClient($addr, $opts);
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Get token
		$user = new Models\Secauthn\User_3();
		$user->setUsername($serverusername);
		$user->setPassword($serverpassword);

		$request = new Services\Secauthn\LogonUserRequest();
		$request->setUser($user);

		list($response, $status) = $secauthnClient->LogonUser($request)->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			throw new Exception("Unable to get token: " . $status->code . ", " . $status->details);
		}

		$token = $response->getToken();
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Get user local groups (and compare with the product's local groups)
		$user = new Models\Netmgmt\User_1();
		$user->setUsername($username);

		$request = new Services\Netmgmt\GetUserLocalGroupsRequest();
		$request->setUser($user);

		list($response, $status) = $netmgmtClient->GetUserLocalGroups($request, ["authorization" => ["Bearer " . $token]])->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			throw new Exception("ERROR [GetUserLocalGroups]: ({$status->code}) {$status->details}");
		}

		$userLocalGroupsNames = [];

		$userLocalGroups = $response->getLocalGroups();
		foreach ($userLocalGroups as $userLocalGroup) {
			array_push($userLocalGroupsNames, $userLocalGroup->getGroupname());
		}

		$localGroupsNamesToRemoveFrom = array_udiff($userLocalGroupsNames, $localGroupsNames, "strcasecmp");
		$localGroupsNamesToAddTo = array_udiff($localGroupsNames, $userLocalGroupsNames, "strcasecmp");

		throw new Exception(json_encode([$localGroupsNamesToRemoveFrom, $localGroupsNamesToAddTo]));
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Get system volumes paths (and compare with the product's volumes paths)
		$request = new Services\Fileio\GetVolumesRequest();

		list($response, $status) = $fileioClient->GetVolumes($request, ["authorization" => ["Bearer " . $token]])->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			throw new Exception("ERROR [GetVolumes]: ({$status->code}) {$status->details}");
		}

		$quotaVolumesToDelete = [];

		$volumes = $response->getVolumes();
		foreach ($volumes as $volume) {
			$volumePaths = $volume->getPaths();
			if (count($volumePaths) === 0) {
				continue;
			}

			foreach ($volumePaths as $volumePath) {
				if (count(array_uintersect([$volumePath], $quotaVolumes, "strcasecmp")) > 0) {
					continue 2; // continue the outer loop => get next $volume
				}
			}

			array_push($quotaVolumesToDelete, $volumePaths[0]);
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Delete non-qualified user quota entries
		foreach ($quotaVolumesToDelete as $volumePath) {
			$user = new Models\Fileio\User_1();
			$user->setUsername($username);

			$request = new Services\Fileio\DeleteUserQuotaEntryRequest();
			$request->setVolumePath($volumePath);
			$request->setUser($user);

			list($response, $status) = $fileioClient->DeleteUserQuotaEntry($request, ["authorization" => ["Bearer " . $token]])->wait();
			if ($status->code !== Grpc\STATUS_OK) {
				throw new Exception("Partial Upgrade/Downgrade occurred. ERROR [DeleteUserQuotaEntry] [{$volumePath}]: ({$status->code}) {$status->details}");
			}
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Remove user from non-qualified local groups
		foreach ($localGroupsNamesToRemoveFrom as $localGroupName) {
			$groupname = $localGroupName;

			$user = new Models\Netmgmt\User_1();
			$user->setUsername($username);

			$localGroup = new Models\Netmgmt\LocalGroup_1();
			$localGroup->setGroupname($groupname);

			$request = new Services\Netmgmt\RemoveUserFromLocalGroupRequest();
			$request->setUser($user);
			$request->setLocalGroup($localGroup);

			list($response, $status) = $netmgmtClient->RemoveUserFromLocalGroup($request, ["authorization" => ["Bearer " . $token]])->wait();
			if ($status->code !== Grpc\STATUS_OK) {
				throw new Exception("Partial Upgrade/Downgrade occurred. ERROR [RemoveUserFromLocalGroup] [{$groupname}]: ({$status->code}) {$status->details}");
			}
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Add user to missing local groups
		foreach ($localGroupsNamesToAddTo as $localGroupName) {
			$groupname = $localGroupName;

			$user = new Models\Netmgmt\User_1();
			$user->setUsername($username);

			$localGroup = new Models\Netmgmt\LocalGroup_1();
			$localGroup->setGroupname($groupname);

			$request = new Services\Netmgmt\AddUserToLocalGroupRequest();
			$request->setUser($user);
			$request->setLocalGroup($localGroup);

			list($response, $status) = $netmgmtClient->AddUserToLocalGroup($request, ["authorization" => ["Bearer " . $token]])->wait();
			if ($status->code !== Grpc\STATUS_OK) {
				throw new Exception("Partial Upgrade/Downgrade occurred. ERROR [AddUserToLocalGroup] [{$groupname}]: ({$status->code}) {$status->details}");
			}
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Set user quota entries
		foreach ($quotaVolumes as $key => $volumePath) {
			$user = new Models\Fileio\User_1();
			$user->setUsername($username);

			$quotaEntry = new Models\Fileio\QuotaEntry_6();
			$quotaEntry->setQuotaThreshold($quotaVolumesThresholds[$key]);
			$quotaEntry->setQuotaLimit($quotaVolumesLimits[$key]);

			$request = new Services\Fileio\SetUserQuotaEntryRequest();
			$request->setVolumePath($volumePath);
			$request->setUser($user);
			$request->setQuotaEntry($quotaEntry);

			list($response, $status) = $fileioClient->SetUserQuotaEntry($request, ["authorization" => ["Bearer " . $token]])->wait();
			if ($status->code !== Grpc\STATUS_OK) {
				throw new Exception("Partial Upgrade/Downgrade occurred. ERROR [SetUserQuotaEntry] [{$volumePath}]: ({$status->code}) {$status->details}");
			}
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	} catch (Exception $e) {
		// Record the error in WHMCS's module log.
		logModuleCall(
			"rdpcloud",
			__FUNCTION__,
			$params,
			$e->getMessage(),
			$e->getTraceAsString()
		);
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		return $e->getMessage();
	}

	return "success";
}

function rdpcloud_ChangePassword(array $params) {
	try {
		// Check server config
		if ($params["server"] !== true) {
			throw new Exception("No server is assigned");
		}

		if ($params["serversecure"] !== true) {
			throw new Exception("SSL/TLS connection must be enabled in the Server Configuration");
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Parse parameters
		$addr = $params["serverip"] . ":" . $params["serverport"];

		$serverusername = $params["serverusername"];
		$serverpassword = $params["serverpassword"];

		$username = $params["username"];
		$password = $params["password"];
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Create clients
		$opts = getOpts();

		$secauthnClient = new Services\Secauthn\SecauthnClient($addr, $opts);
		$netmgmtClient = new Services\Netmgmt\NetmgmtClient($addr, $opts);
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Get token
		$user = new Models\Secauthn\User_3();
		$user->setUsername($serverusername);
		$user->setPassword($serverpassword);

		$request = new Services\Secauthn\LogonUserRequest();
		$request->setUser($user);

		list($response, $status) = $secauthnClient->LogonUser($request)->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			throw new Exception("Unable to get token: " . $status->code . ", " . $status->details);
		}

		$token = $response->getToken();
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Change user password
		$user = new Models\Netmgmt\User_3();
		$user->setUsername($username);
		$user->setPassword($password);

		$request = new Services\Netmgmt\ChangeUserPasswordRequest();
		$request->setUser($user);

		list($response, $status) = $netmgmtClient->ChangeUserPassword($request, ["authorization" => ["Bearer " . $token]])->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			throw new Exception("ERROR [ChangeUserPassword]: " . $status->code . ", " . $status->details);
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	} catch (Exception $e) {
		// Record the error in WHMCS's module log.
		logModuleCall(
			"rdpcloud",
			__FUNCTION__,
			$params,
			$e->getMessage(),
			$e->getTraceAsString()
		);
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		return $e->getMessage();
	}

	return "success";
}

function rdpcloud_SuspendAccount($params) {
	return "success";
}

function rdpcloud_UnsuspendAccount($params) {
	return "success";
}
