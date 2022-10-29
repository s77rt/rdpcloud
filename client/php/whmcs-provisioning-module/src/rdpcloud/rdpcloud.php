<?php

if (!defined("WHMCS")) {
	die("This file cannot be accessed directly");
}

require dirname(__FILE__) . '/vendor/autoload.php';

function getOpts() {
	return [
		"credentials" => Grpc\ChannelCredentials::createSsl(file_get_contents(dirname(__FILE__) . '/cert/server-cert.pem')),
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
	$localGroups = array_filter(array_map("trim", preg_split('/\r\n|\r|\n/', $params["configoption1"])));
	$localGroups__Groupname = $localGroups;

	$diskQuotas = array_filter(array_map("trim", preg_split('/\r\n|\r|\n/', $params["configoption2"])));
	$diskQuotas__VolumePath = [];
	$diskQuotas__QuotaThreshold = [];
	$diskQuotas__QuotaLimit = [];

	foreach ($diskQuotas as $diskQuota) {
		$parsedDiskQuota = array_filter(array_map("trim", explode(",", $diskQuota)));
		if (count($parsedDiskQuota) != 3) {
			return [null, null, null, null, "Disk Quotas must be in the form VolumePath, Quota Threshold, Quota Limit"];
		}

		$volumePath = $parsedDiskQuota[0];
		$quotaThreshold = parseSize($parsedDiskQuota[1]);
		$quotaLimit = parseSize($parsedDiskQuota[2]);

		if ($quotaThreshold === 0 || $quotaLimit === 0) {
			return [null, null, null, null, "Disk Quotas contains invalid Quota Threshold and/or Quota Limit values"];
		}

		array_push($diskQuotas__VolumePath, $volumePath);
		array_push($diskQuotas__QuotaThreshold, $quotaThreshold);
		array_push($diskQuotas__QuotaLimit, $quotaLimit);
	}

	return [$localGroups__Groupname, $diskQuotas__VolumePath, $diskQuotas__QuotaThreshold, $diskQuotas__QuotaLimit, null];
}

// genUsernameAndPassword generates a username and a password that meets the Windows password policy requirements
// a new username will be generared only if $username_level is greater than zero
function genUsernameAndPassword(array $params, int $username_level = 0) {
	$username = $params["username"];
	$password = $params["password"];

	if ($username_level > 0) {
		$firstname = ucfirst(preg_replace("/[^a-zA-Z0-9]+/", "", $params["clientsdetails"]["firstname"]));
		$lastname = ucfirst(preg_replace("/[^a-zA-Z0-9]+/", "", $params["clientsdetails"]["lastname"]));
		$emailname = "";
		if (filter_var($params["clientsdetails"]["email"], FILTER_VALIDATE_EMAIL) !== false) {
			$email = $params["clientsdetails"]["email"];
			$emailname = substr($email, 0, strrpos($email, "@"));
			$emailname = ucfirst(preg_replace("/[^a-zA-Z0-9]+/", "", $emailname));
		}

		// Username min and max length without the extra digits
		$USERNAME_MIN_LENGTH = 7;
		$USERNAME_MAX_LENGTH = 15;

		$USERNAME_EXTRA_DIGITS = 2;

		$username = "";

		if (strlen($username) < $USERNAME_MIN_LENGTH || $username_level >= 1) {
			$username .= $firstname;
		}
		if (strlen($username) < $USERNAME_MIN_LENGTH || $username_level >= 2) {
			$username .= $lastname;
		}
		if (strlen($username) < $USERNAME_MIN_LENGTH || $username_level >= 3) {
			$username .= $emailname;
		}
		if (strlen($username) < $USERNAME_MIN_LENGTH || $username_level >= 4) {
			$username = "User";
		}

		// Min Length Check
		$missing_length = $USERNAME_MIN_LENGTH - strlen($username);
		if ($missing_length > 0) {
			$username .= random_int(pow(10, $missing_length - 1), pow(10, $missing_length) - 1);
		}

		// Max Length Check
		$username = substr($username, 0, $USERNAME_MAX_LENGTH);

		// Append extra digits
		$username .= random_int(pow(10, $USERNAME_EXTRA_DIGITS - 1), pow(10, $USERNAME_EXTRA_DIGITS) - 1);
	}

	$old_password = $params["password"];

	$KEYSPACE_UPPER_CASE = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
	$KEYSPACE_UPPER_CASE__LENGTH = strlen($KEYSPACE_UPPER_CASE);
	$KEYSPACE_LOWER_CASE = 'abcdefghijklmnopqrstuvwxyz';
	$KEYSPACE_LOWER_CASE__LENGTH = strlen($KEYSPACE_LOWER_CASE);
	$KEYSPACE_NUMBERS = '0123456789';
	$KEYSPACE_NUMBERS__LENGTH = strlen($KEYSPACE_NUMBERS);
	$KEYSPACE_SYMBOLS = '~!@#$%^*_-+=`|(){}[]:;,.?';
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
		'AdminSingleSignOnLabel' => 'Login to RDP Control Panel',
		'ServiceSingleSignOnLabel' => 'Login to RDP Control Panel',
	);
}

function rdpcloud_ConfigOptions() {
	return array(
		'Local Groups' => array(
			'FriendlyName' => 'Local Groups<br>(One per line)',
			'Type' => 'textarea',
			'Description' => 'Local Groups',
			'Default' => <<<'EOD'
Users
Remote Desktop Users
EOD,
		),
		'Disk Quotas'=>array(
			'FriendlyName'=>'Disk Quotas<br>(One per line)',
			'Type'=>'textarea',
			'Description'=>'Format: VolumePath, Quota Threshold, Quota Limit',
			'Default'=>'C:\\, 1.5GB, 2GB',
		),
		'Enable Single Sign-On'=>array(
			'FriendlyName'=>'Enable Single Sign-On',
			'Type'=>'yesno',
			'Description'=>'Enable Single Sign-On',
			'Default'=>true,
		),
		'RDP Control Panel URL'=>array(
			'FriendlyName'=>'RDP Control Panel URL',
			'Type'=>'text',
		),
	);
}

function rdpcloud_AdminSingleSignOn(array $params) {
	if ($params["configoption3"] !== true && $params["configoption3"] !== "on") {
		return array(
			'success' => false,
			'errorMsg' => "Single Sign-On is not enabled",
		);
	}
	if (filter_var($params["configoption4"], FILTER_VALIDATE_URL) === FALSE) {
		return array(
			'success' => false,
			'errorMsg' => "RDP Control Panel URL is not valid",
		);
	}

	return array(
		'success' => true,
		'redirectTo' => $params["configoption4"],
	);
}

function rdpcloud_ServiceSingleSignOn(array $params) {
	if ($params["configoption3"] !== true && $params["configoption3"] !== "on") {
		return array(
			'success' => false,
			'errorMsg' => "Single Sign-On is not enabled",
		);
	}
	if (filter_var($params["configoption4"], FILTER_VALIDATE_URL) === FALSE) {
		return array(
			'success' => false,
			'errorMsg' => "RDP Control Panel URL is not valid",
		);
	}

	return array(
		'success' => true,
		'redirectTo' => $params["configoption4"],
	);
}

function rdpcloud_AdminCustomButtonArray() {
	$buttonArray = array(
		"Initiate Server Restart" => "InitiateServerRestart",
		"Abort Server Restart" => "AbortServerRestart",
	);
	return $buttonArray;
}

function rdpcloud_ClientAreaCustomButtonArray() {
	$buttonArray = array(
		"Login to RDP Control Panel" => "LoginRDPControlPanel",
	);
	return $buttonArray;
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

		list($localGroups__Groupname, $diskQuotas__VolumePath, $diskQuotas__QuotaThreshold, $diskQuotas__QuotaLimit, $error) = parseConfigOptions($params);
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
			throw new Exception("[LogonUser]: ({$status->code}) {$status->details}");
		}

		$token = $response->getToken();
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Add user
		$username_level = 0;
		if (strlen($username) === 0) {
			$username_level = 1;
		}
		while (true) {
			list($username, $password) = genUsernameAndPassword($params, $username_level);
			$params["model"]->serviceProperties->save(["Username" => $username, "Password" => $password]);

			$user = new Models\Netmgmt\User_3();
			$user->setUsername($username);
			$user->setPassword($password);

			$request = new Services\Netmgmt\AddUserRequest();
			$request->setUser($user);

			list($response, $status) = $netmgmtClient->AddUser($request, ["authorization" => ["Bearer " . $token]])->wait();
			if ($status->code !== Grpc\STATUS_OK) {
				if ($status->code === Grpc\STATUS_ALREADY_EXISTS && ($status->details === "User already exists" || $status->details === "Group already exists") && $username_level < 5) {
					$username_level++;
					continue;
				} else {
					throw new Exception("[AddUser]: ({$status->code}) {$status->details}");
				}
			}

			break;
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Add user to local groups
		foreach ($localGroups__Groupname as $localGroup__Groupname) {
			$groupname = $localGroup__Groupname;

			$user = new Models\Netmgmt\User_1();
			$user->setUsername($username);

			$localGroup = new Models\Netmgmt\LocalGroup_1();
			$localGroup->setGroupname($groupname);

			$request = new Services\Netmgmt\AddUserToLocalGroupRequest();
			$request->setUser($user);
			$request->setLocalGroup($localGroup);

			list($response, $status) = $netmgmtClient->AddUserToLocalGroup($request, ["authorization" => ["Bearer " . $token]])->wait();
			if ($status->code !== Grpc\STATUS_OK) {
				throw new Exception("User created but an error occurred. [AddUserToLocalGroup] [{$groupname}]: ({$status->code}) {$status->details}. IMPORTANT: DELETE THE USER BEFORE TRYING AGAIN");
			}
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Set user quota entries
		foreach ($diskQuotas__VolumePath as $key => $diskQuota__VolumePath) {
			$volumePath = $diskQuota__VolumePath;
			$quotaThreshold = $diskQuotas__QuotaThreshold[$key];
			$quotaLimit = $diskQuotas__QuotaLimit[$key];

			$user = new Models\Fileio\User_1();
			$user->setUsername($username);

			$quotaEntry = new Models\Fileio\QuotaEntry_6();
			$quotaEntry->setQuotaThreshold($quotaThreshold);
			$quotaEntry->setQuotaLimit($quotaLimit);

			$request = new Services\Fileio\SetUserQuotaEntryRequest();
			$request->setVolumePath($volumePath);
			$request->setUser($user);
			$request->setQuotaEntry($quotaEntry);

			list($response, $status) = $fileioClient->SetUserQuotaEntry($request, ["authorization" => ["Bearer " . $token]])->wait();
			if ($status->code !== Grpc\STATUS_OK) {
				throw new Exception("User created but an error occurred. [SetUserQuotaEntry] [{$volumePath}]: ({$status->code}) {$status->details}. IMPORTANT: DELETE THE USER BEFORE TRYING AGAIN");
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
			throw new Exception("[LogonUser]: ({$status->code}) {$status->details}");
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
				throw new Exception("[DisableUser]: ({$status->code}) {$status->details}");
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
				throw new Exception("[LogoffUser]: ({$status->code}) {$status->details}");
			}
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Get user sid
		$request = new Services\Secauthz\LookupAccountSidByUsernameRequest();
		$request->setUsername($username);

		list($response, $status) = $secauthzClient->LookupAccountSidByUsername($request, ["authorization" => ["Bearer " . $token]])->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			throw new Exception("[LookupAccountSidByUsername]: ({$status->code}) {$status->details}");
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
			throw new Exception("[GetVolumes]: ({$status->code}) {$status->details}");
		}

		$volumesPaths = [];

		$volumes = $response->getVolumes();
		foreach ($volumes as $volume) {
			$volumePaths = $volume->getPaths();
			if (count($volumePaths) === 0) {
				continue;
			}

			array_push($volumesPaths, $volumePaths[0]);
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
				throw new Exception("[DeleteUserQuotaEntry] [{$volumePath}]: ({$status->code}) {$status->details}");
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
			throw new Exception("[DeleteUser]: ({$status->code}) {$status->details}");
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

		list($localGroups__Groupname, $diskQuotas__VolumePath, $diskQuotas__QuotaThreshold, $diskQuotas__QuotaLimit, $error) = parseConfigOptions($params);
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
			throw new Exception("[LogonUser]: ({$status->code}) {$status->details}");
		}

		$token = $response->getToken();
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Get user local groups (and compare with the product's)
		$user = new Models\Netmgmt\User_1();
		$user->setUsername($username);

		$request = new Services\Netmgmt\GetUserLocalGroupsRequest();
		$request->setUser($user);

		list($response, $status) = $netmgmtClient->GetUserLocalGroups($request, ["authorization" => ["Bearer " . $token]])->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			throw new Exception("[GetUserLocalGroups]: ({$status->code}) {$status->details}");
		}

		$userLocalGroups__Groupname = [];

		$userLocalGroups = $response->getLocalGroups();
		foreach ($userLocalGroups as $userLocalGroup) {
			array_push($userLocalGroups__Groupname, $userLocalGroup->getGroupname());
		}

		$localGroupsToRemoveFrom__Groupname = array_udiff($userLocalGroups__Groupname, $localGroups__Groupname, "strcasecmp");
		$localGroupsToAddTo__Groupname = array_udiff($localGroups__Groupname, $userLocalGroups__Groupname, "strcasecmp");

		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Get system volumes (and compare with the product's)
		$request = new Services\Fileio\GetVolumesRequest();

		list($response, $status) = $fileioClient->GetVolumes($request, ["authorization" => ["Bearer " . $token]])->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			throw new Exception("[GetVolumes]: ({$status->code}) {$status->details}");
		}

		$diskQuotasToDelete__VolumePath = [];

		$volumes = $response->getVolumes();
		foreach ($volumes as $volume) {
			$volumePaths = $volume->getPaths();
			if (count($volumePaths) === 0) {
				continue;
			}

			if (count(array_uintersect($volumePaths, $diskQuotas__VolumePath, "strcasecmp")) > 0) {
				continue;
			}

			array_push($diskQuotasToDelete__VolumePath, $volumePaths[0]);
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Delete non-qualified user quota entries
		foreach ($diskQuotasToDelete__VolumePath as $diskQuotaToDelete__VolumePath) {
			$volumePath = $diskQuotaToDelete__VolumePath;

			$user = new Models\Fileio\User_1();
			$user->setUsername($username);

			$request = new Services\Fileio\DeleteUserQuotaEntryRequest();
			$request->setVolumePath($volumePath);
			$request->setUser($user);

			list($response, $status) = $fileioClient->DeleteUserQuotaEntry($request, ["authorization" => ["Bearer " . $token]])->wait();
			if ($status->code !== Grpc\STATUS_OK) {
				throw new Exception("Partial Upgrade/Downgrade occurred. [DeleteUserQuotaEntry] [{$volumePath}]: ({$status->code}) {$status->details}");
			}
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Remove user from non-qualified local groups
		foreach ($localGroupsToRemoveFrom__Groupname as $localGroupToRemoveFrom__Groupname) {
			$groupname = $localGroupToRemoveFrom__Groupname;

			$user = new Models\Netmgmt\User_1();
			$user->setUsername($username);

			$localGroup = new Models\Netmgmt\LocalGroup_1();
			$localGroup->setGroupname($groupname);

			$request = new Services\Netmgmt\RemoveUserFromLocalGroupRequest();
			$request->setUser($user);
			$request->setLocalGroup($localGroup);

			list($response, $status) = $netmgmtClient->RemoveUserFromLocalGroup($request, ["authorization" => ["Bearer " . $token]])->wait();
			if ($status->code !== Grpc\STATUS_OK) {
				throw new Exception("Partial Upgrade/Downgrade occurred. [RemoveUserFromLocalGroup] [{$groupname}]: ({$status->code}) {$status->details}");
			}
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Add user to missing local groups
		foreach ($localGroupsToAddTo__Groupname as $localGroupToAddTo__Groupname) {
			$groupname = $localGroupToAddTo__Groupname;

			$user = new Models\Netmgmt\User_1();
			$user->setUsername($username);

			$localGroup = new Models\Netmgmt\LocalGroup_1();
			$localGroup->setGroupname($groupname);

			$request = new Services\Netmgmt\AddUserToLocalGroupRequest();
			$request->setUser($user);
			$request->setLocalGroup($localGroup);

			list($response, $status) = $netmgmtClient->AddUserToLocalGroup($request, ["authorization" => ["Bearer " . $token]])->wait();
			if ($status->code !== Grpc\STATUS_OK) {
				throw new Exception("Partial Upgrade/Downgrade occurred. [AddUserToLocalGroup] [{$groupname}]: ({$status->code}) {$status->details}");
			}
		}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Set user quota entries
		foreach ($diskQuotas__VolumePath as $key => $diskQuota__VolumePath) {
			$volumePath = $diskQuota__VolumePath;
			$quotaThreshold = $diskQuotas__QuotaThreshold[$key];
			$quotaLimit = $diskQuotas__QuotaLimit[$key];

			$user = new Models\Fileio\User_1();
			$user->setUsername($username);

			$quotaEntry = new Models\Fileio\QuotaEntry_6();
			$quotaEntry->setQuotaThreshold($quotaThreshold);
			$quotaEntry->setQuotaLimit($quotaLimit);

			$request = new Services\Fileio\SetUserQuotaEntryRequest();
			$request->setVolumePath($volumePath);
			$request->setUser($user);
			$request->setQuotaEntry($quotaEntry);

			list($response, $status) = $fileioClient->SetUserQuotaEntry($request, ["authorization" => ["Bearer " . $token]])->wait();
			if ($status->code !== Grpc\STATUS_OK) {
				throw new Exception("Partial Upgrade/Downgrade occurred. [SetUserQuotaEntry] [{$volumePath}]: ({$status->code}) {$status->details}");
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
			throw new Exception("[LogonUser]: ({$status->code}) {$status->details}");
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
			throw new Exception("[ChangeUserPassword]: ({$status->code}) {$status->details}");
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
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Get token
		$user = new Models\Secauthn\User_3();
		$user->setUsername($serverusername);
		$user->setPassword($serverpassword);

		$request = new Services\Secauthn\LogonUserRequest();
		$request->setUser($user);

		list($response, $status) = $secauthnClient->LogonUser($request)->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			throw new Exception("[LogonUser]: ({$status->code}) {$status->details}");
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
				throw new Exception("[DisableUser]: ({$status->code}) {$status->details}");
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
				throw new Exception("[LogoffUser]: ({$status->code}) {$status->details}");
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

function rdpcloud_UnsuspendAccount($params) {
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
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Get token
		$user = new Models\Secauthn\User_3();
		$user->setUsername($serverusername);
		$user->setPassword($serverpassword);

		$request = new Services\Secauthn\LogonUserRequest();
		$request->setUser($user);

		list($response, $status) = $secauthnClient->LogonUser($request)->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			throw new Exception("[LogonUser]: ({$status->code}) {$status->details}");
		}

		$token = $response->getToken();
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Enable user
		$user = new Models\Netmgmt\User_1();
		$user->setUsername($username);

		$request = new Services\Netmgmt\EnableUserRequest();
		$request->setUser($user);

		list($response, $status) = $netmgmtClient->EnableUser($request, ["authorization" => ["Bearer " . $token]])->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			if ($status->code === Grpc\STATUS_FAILED_PRECONDITION && $status->details === "User is already enabled") {
				;
			} else {
				throw new Exception("[EnableUser]: ({$status->code}) {$status->details}");
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

function rdpcloud_InitiateServerRestart($params) {
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
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Create clients
		$opts = getOpts();

		$secauthnClient = new Services\Secauthn\SecauthnClient($addr, $opts);
		$shutdownClient = new Services\Shutdown\ShutdownClient($addr, $opts);
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Get token
		$user = new Models\Secauthn\User_3();
		$user->setUsername($serverusername);
		$user->setPassword($serverpassword);

		$request = new Services\Secauthn\LogonUserRequest();
		$request->setUser($user);

		list($response, $status) = $secauthnClient->LogonUser($request)->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			throw new Exception("[LogonUser]: ({$status->code}) {$status->details}");
		}

		$token = $response->getToken();
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Initiate system shutdown
		$request = new Services\Shutdown\InitiateSystemShutdownRequest();
		$request->setMessage("Restarting the server in 90 seconds. Sorry for any inconvenience this may cause");
		$request->setTimeout(90);
		$request->setForce(true);
		$request->setReboot(true);
		$request->setReason(0x00050000);

		list($response, $status) = $shutdownClient->InitiateSystemShutdown($request, ["authorization" => ["Bearer " . $token]])->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			throw new Exception("[InitiateSystemShutdown]: ({$status->code}) {$status->details}");
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

function rdpcloud_AbortServerRestart($params) {
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
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Create clients
		$opts = getOpts();

		$secauthnClient = new Services\Secauthn\SecauthnClient($addr, $opts);
		$shutdownClient = new Services\Shutdown\ShutdownClient($addr, $opts);
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Get token
		$user = new Models\Secauthn\User_3();
		$user->setUsername($serverusername);
		$user->setPassword($serverpassword);

		$request = new Services\Secauthn\LogonUserRequest();
		$request->setUser($user);

		list($response, $status) = $secauthnClient->LogonUser($request)->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			throw new Exception("[LogonUser]: ({$status->code}) {$status->details}");
		}

		$token = $response->getToken();
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		// Abort system shutdown
		$request = new Services\Shutdown\AbortSystemShutdownRequest();

		list($response, $status) = $shutdownClient->AbortSystemShutdown($request, ["authorization" => ["Bearer " . $token]])->wait();
		if ($status->code !== Grpc\STATUS_OK) {
			throw new Exception("[AbortSystemShutdown]: ({$status->code}) {$status->details}");
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

function rdpcloud_LoginRDPControlPanel(array $params) {
	$rdpControlPanelURL = null;
	if ($params["configoption3"] === true || $params["configoption3"] === "on") {
		$rdpControlPanelURL = filter_var($params["configoption4"], FILTER_VALIDATE_URL);
	}

	$templateFile = "templates/rdp-control-panel";
	return array(
		'templatefile' => $templateFile,
		'vars' => array(
			'rdpControlPanelURL' => $rdpControlPanelURL,
		),
	);
}
