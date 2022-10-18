<?php
// GENERATED CODE -- DO NOT EDIT!

namespace Services\Shutdown;

/**
 */
class ShutdownClient extends \Grpc\BaseStub {

    /**
     * @param string $hostname hostname
     * @param array $opts channel options
     * @param \Grpc\Channel $channel (optional) re-use channel object
     */
    public function __construct($hostname, $opts, $channel = null) {
        parent::__construct($hostname, $opts, $channel);
    }

    /**
     * @param \Services\Shutdown\InitiateSystemShutdownRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function InitiateSystemShutdown(\Services\Shutdown\InitiateSystemShutdownRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.shutdown.Shutdown/InitiateSystemShutdown',
        $argument,
        ['\Services\Shutdown\InitiateSystemShutdownResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Shutdown\AbortSystemShutdownRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function AbortSystemShutdown(\Services\Shutdown\AbortSystemShutdownRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.shutdown.Shutdown/AbortSystemShutdown',
        $argument,
        ['\Services\Shutdown\AbortSystemShutdownResponse', 'decode'],
        $metadata, $options);
    }

}
