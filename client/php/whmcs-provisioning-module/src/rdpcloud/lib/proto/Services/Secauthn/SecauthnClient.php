<?php
// GENERATED CODE -- DO NOT EDIT!

namespace Services\Secauthn;

/**
 */
class SecauthnClient extends \Grpc\BaseStub {

    /**
     * @param string $hostname hostname
     * @param array $opts channel options
     * @param \Grpc\Channel $channel (optional) re-use channel object
     */
    public function __construct($hostname, $opts, $channel = null) {
        parent::__construct($hostname, $opts, $channel);
    }

    /**
     * @param \Services\Secauthn\LogonUserRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function LogonUser(\Services\Secauthn\LogonUserRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.secauthn.Secauthn/LogonUser',
        $argument,
        ['\Services\Secauthn\LogonUserResponse', 'decode'],
        $metadata, $options);
    }

}
