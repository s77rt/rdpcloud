<?php
// GENERATED CODE -- DO NOT EDIT!

namespace Services\Fileio;

/**
 */
class FileioClient extends \Grpc\BaseStub {

    /**
     * @param string $hostname hostname
     * @param array $opts channel options
     * @param \Grpc\Channel $channel (optional) re-use channel object
     */
    public function __construct($hostname, $opts, $channel = null) {
        parent::__construct($hostname, $opts, $channel);
    }

    /**
     * @param \Services\Fileio\GetQuotaStateRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function GetQuotaState(\Services\Fileio\GetQuotaStateRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.fileio.Fileio/GetQuotaState',
        $argument,
        ['\Services\Fileio\GetQuotaStateResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Fileio\SetQuotaStateRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function SetQuotaState(\Services\Fileio\SetQuotaStateRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.fileio.Fileio/SetQuotaState',
        $argument,
        ['\Services\Fileio\SetQuotaStateResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Fileio\GetDefaultQuotaRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function GetDefaultQuota(\Services\Fileio\GetDefaultQuotaRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.fileio.Fileio/GetDefaultQuota',
        $argument,
        ['\Services\Fileio\GetDefaultQuotaResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Fileio\SetDefaultQuotaRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function SetDefaultQuota(\Services\Fileio\SetDefaultQuotaRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.fileio.Fileio/SetDefaultQuota',
        $argument,
        ['\Services\Fileio\SetDefaultQuotaResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Fileio\GetUsersQuotaEntriesRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function GetUsersQuotaEntries(\Services\Fileio\GetUsersQuotaEntriesRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.fileio.Fileio/GetUsersQuotaEntries',
        $argument,
        ['\Services\Fileio\GetUsersQuotaEntriesResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Fileio\GetUserQuotaEntryRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function GetUserQuotaEntry(\Services\Fileio\GetUserQuotaEntryRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.fileio.Fileio/GetUserQuotaEntry',
        $argument,
        ['\Services\Fileio\GetUserQuotaEntryResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Fileio\GetMyUserQuotaEntryRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function GetMyUserQuotaEntry(\Services\Fileio\GetMyUserQuotaEntryRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.fileio.Fileio/GetMyUserQuotaEntry',
        $argument,
        ['\Services\Fileio\GetMyUserQuotaEntryResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Fileio\SetUserQuotaEntryRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function SetUserQuotaEntry(\Services\Fileio\SetUserQuotaEntryRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.fileio.Fileio/SetUserQuotaEntry',
        $argument,
        ['\Services\Fileio\SetUserQuotaEntryResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Fileio\DeleteUserQuotaEntryRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function DeleteUserQuotaEntry(\Services\Fileio\DeleteUserQuotaEntryRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.fileio.Fileio/DeleteUserQuotaEntry',
        $argument,
        ['\Services\Fileio\DeleteUserQuotaEntryResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Fileio\GetVolumesRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function GetVolumes(\Services\Fileio\GetVolumesRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.fileio.Fileio/GetVolumes',
        $argument,
        ['\Services\Fileio\GetVolumesResponse', 'decode'],
        $metadata, $options);
    }

}
