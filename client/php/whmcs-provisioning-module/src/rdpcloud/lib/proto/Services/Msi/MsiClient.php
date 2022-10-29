<?php
// GENERATED CODE -- DO NOT EDIT!

namespace Services\Msi;

/**
 */
class MsiClient extends \Grpc\BaseStub {

    /**
     * @param string $hostname hostname
     * @param array $opts channel options
     * @param \Grpc\Channel $channel (optional) re-use channel object
     */
    public function __construct($hostname, $opts, $channel = null) {
        parent::__construct($hostname, $opts, $channel);
    }

    /**
     * @param \Services\Msi\GetProductsRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function GetProducts(\Services\Msi\GetProductsRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/services.msi.Msi/GetProducts',
        $argument,
        ['\Services\Msi\GetProductsResponse', 'decode'],
        $metadata, $options);
    }

}
