<?php
// GENERATED CODE -- DO NOT EDIT!

namespace Services\Sysinfo;

/**
 */
class SysinfoClient extends \Grpc\BaseStub {

    /**
     * @param string $hostname hostname
     * @param array $opts channel options
     * @param \Grpc\Channel $channel (optional) re-use channel object
     */
    public function __construct($hostname, $opts, $channel = null) {
        parent::__construct($hostname, $opts, $channel);
    }

    /**
     * @param \Services\Sysinfo\GetUptimeRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function GetUptime(\Services\Sysinfo\GetUptimeRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.sysinfo.Sysinfo/GetUptime',
        $argument,
        ['\Services\Sysinfo\GetUptimeResponse', 'decode'],
        $metadata, $options);
    }

}
