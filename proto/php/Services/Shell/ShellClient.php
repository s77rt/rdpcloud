<?php
// GENERATED CODE -- DO NOT EDIT!

namespace Services\Shell;

/**
 */
class ShellClient extends \Grpc\BaseStub {

    /**
     * @param string $hostname hostname
     * @param array $opts channel options
     * @param \Grpc\Channel $channel (optional) re-use channel object
     */
    public function __construct($hostname, $opts, $channel = null) {
        parent::__construct($hostname, $opts, $channel);
    }

    /**
     * @param \Services\Shell\DeleteProfileRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function DeleteProfile(\Services\Shell\DeleteProfileRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.shell.Shell/DeleteProfile',
        $argument,
        ['\Services\Shell\DeleteProfileResponse', 'decode'],
        $metadata, $options);
    }

}
