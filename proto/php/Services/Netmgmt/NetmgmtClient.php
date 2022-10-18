<?php
// GENERATED CODE -- DO NOT EDIT!

namespace Services\Netmgmt;

/**
 */
class NetmgmtClient extends \Grpc\BaseStub {

    /**
     * @param string $hostname hostname
     * @param array $opts channel options
     * @param \Grpc\Channel $channel (optional) re-use channel object
     */
    public function __construct($hostname, $opts, $channel = null) {
        parent::__construct($hostname, $opts, $channel);
    }

    /**
     * @param \Services\Netmgmt\AddUserRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function AddUser(\Services\Netmgmt\AddUserRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.netmgmt.Netmgmt/AddUser',
        $argument,
        ['\Services\Netmgmt\AddUserResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Netmgmt\DeleteUserRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function DeleteUser(\Services\Netmgmt\DeleteUserRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.netmgmt.Netmgmt/DeleteUser',
        $argument,
        ['\Services\Netmgmt\DeleteUserResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Netmgmt\GetUsersRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function GetUsers(\Services\Netmgmt\GetUsersRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.netmgmt.Netmgmt/GetUsers',
        $argument,
        ['\Services\Netmgmt\GetUsersResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Netmgmt\GetUserRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function GetUser(\Services\Netmgmt\GetUserRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.netmgmt.Netmgmt/GetUser',
        $argument,
        ['\Services\Netmgmt\GetUserResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Netmgmt\GetMyUserRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function GetMyUser(\Services\Netmgmt\GetMyUserRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.netmgmt.Netmgmt/GetMyUser',
        $argument,
        ['\Services\Netmgmt\GetMyUserResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Netmgmt\GetUserLocalGroupsRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function GetUserLocalGroups(\Services\Netmgmt\GetUserLocalGroupsRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.netmgmt.Netmgmt/GetUserLocalGroups',
        $argument,
        ['\Services\Netmgmt\GetUserLocalGroupsResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Netmgmt\GetMyUserLocalGroupsRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function GetMyUserLocalGroups(\Services\Netmgmt\GetMyUserLocalGroupsRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.netmgmt.Netmgmt/GetMyUserLocalGroups',
        $argument,
        ['\Services\Netmgmt\GetMyUserLocalGroupsResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Netmgmt\ChangeUserPasswordRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function ChangeUserPassword(\Services\Netmgmt\ChangeUserPasswordRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.netmgmt.Netmgmt/ChangeUserPassword',
        $argument,
        ['\Services\Netmgmt\ChangeUserPasswordResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Netmgmt\ChangeMyUserPasswordRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function ChangeMyUserPassword(\Services\Netmgmt\ChangeMyUserPasswordRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.netmgmt.Netmgmt/ChangeMyUserPassword',
        $argument,
        ['\Services\Netmgmt\ChangeMyUserPasswordResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Netmgmt\EnableUserRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function EnableUser(\Services\Netmgmt\EnableUserRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.netmgmt.Netmgmt/EnableUser',
        $argument,
        ['\Services\Netmgmt\EnableUserResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Netmgmt\DisableUserRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function DisableUser(\Services\Netmgmt\DisableUserRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.netmgmt.Netmgmt/DisableUser',
        $argument,
        ['\Services\Netmgmt\DisableUserResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Netmgmt\AddUserToLocalGroupRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function AddUserToLocalGroup(\Services\Netmgmt\AddUserToLocalGroupRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.netmgmt.Netmgmt/AddUserToLocalGroup',
        $argument,
        ['\Services\Netmgmt\AddUserToLocalGroupResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Netmgmt\RemoveUserFromLocalGroupRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function RemoveUserFromLocalGroup(\Services\Netmgmt\RemoveUserFromLocalGroupRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.netmgmt.Netmgmt/RemoveUserFromLocalGroup',
        $argument,
        ['\Services\Netmgmt\RemoveUserFromLocalGroupResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Netmgmt\GetLocalGroupsRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function GetLocalGroups(\Services\Netmgmt\GetLocalGroupsRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.netmgmt.Netmgmt/GetLocalGroups',
        $argument,
        ['\Services\Netmgmt\GetLocalGroupsResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Netmgmt\GetUsersInLocalGroupRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function GetUsersInLocalGroup(\Services\Netmgmt\GetUsersInLocalGroupRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.netmgmt.Netmgmt/GetUsersInLocalGroup',
        $argument,
        ['\Services\Netmgmt\GetUsersInLocalGroupResponse', 'decode'],
        $metadata, $options);
    }

}
