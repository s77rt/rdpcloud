<?php
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: services/shutdown/system_shutdown.proto

namespace GPBMetadata\Services\Shutdown;

class SystemShutdown
{
    public static $is_initialized = false;

    public static function initOnce() {
        $pool = \Google\Protobuf\Internal\DescriptorPool::getGeneratedPool();

        if (static::$is_initialized == true) {
          return;
        }
        $pool->internalAddGeneratedFile(
            '
�
\'services/shutdown/system_shutdown.protoservices.shutdown"p
InitiateSystemShutdownRequest
message (	
timeout (
force (
reboot (
reason (
InitiateSystemShutdownResponse"
AbortSystemShutdownRequest"
AbortSystemShutdownResponseB6Z4github.com/s77rt/rdpcloud/proto/go/services/shutdownbproto3'
        , true);

        static::$is_initialized = true;
    }
}
