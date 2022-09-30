<?php
// GENERATED CODE -- DO NOT EDIT!

namespace Services\Secauthz;

/**
 */
class SecauthzClient extends \Grpc\BaseStub {

    /**
     * @param string $hostname hostname
     * @param array $opts channel options
     * @param \Grpc\Channel $channel (optional) re-use channel object
     */
    public function __construct($hostname, $opts, $channel = null) {
        parent::__construct($hostname, $opts, $channel);
    }

    /**
     * @param \Services\Secauthz\LookupAccountSidByUsernameRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function LookupAccountSidByUsername(\Services\Secauthz\LookupAccountSidByUsernameRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.secauthz.Secauthz/LookupAccountSidByUsername',
        $argument,
        ['\Services\Secauthz\LookupAccountSidByUsernameResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Services\Secauthz\LookupAccountUsernameBySidRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function LookupAccountUsernameBySid(\Services\Secauthz\LookupAccountUsernameBySidRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.secauthz.Secauthz/LookupAccountUsernameBySid',
        $argument,
        ['\Services\Secauthz\LookupAccountUsernameBySidResponse', 'decode'],
        $metadata, $options);
    }

}
