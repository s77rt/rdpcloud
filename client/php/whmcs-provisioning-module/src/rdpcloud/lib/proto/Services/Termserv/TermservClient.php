<?php
// GENERATED CODE -- DO NOT EDIT!

namespace Services\Termserv;

/**
 */
class TermservClient extends \Grpc\BaseStub {

    /**
     * @param string $hostname hostname
     * @param array $opts channel options
     * @param \Grpc\Channel $channel (optional) re-use channel object
     */
    public function __construct($hostname, $opts, $channel = null) {
        parent::__construct($hostname, $opts, $channel);
    }

    /**
     * @param \Services\Termserv\LogoffUserRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function LogoffUser(\Services\Termserv\LogoffUserRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.termserv.Termserv/LogoffUser',
        $argument,
        ['\Services\Termserv\LogoffUserResponse', 'decode'],
        $metadata, $options);
    }

}
