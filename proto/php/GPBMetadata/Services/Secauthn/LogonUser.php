<?php
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: services/secauthn/logon_user.proto

namespace GPBMetadata\Services\Secauthn;

class LogonUser
{
    public static $is_initialized = false;

    public static function initOnce() {
        $pool = \Google\Protobuf\Internal\DescriptorPool::getGeneratedPool();

        if (static::$is_initialized == true) {
          return;
        }
        \GPBMetadata\Models\Secauthn\User::initOnce();
        $pool->internalAddGeneratedFile(
            '
�
"services/secauthn/logon_user.protoservices.secauthn"9
LogonUserRequest%
user (2.models.secauthn.User_3""
LogonUserResponse
token (	B6Z4github.com/s77rt/rdpcloud/proto/go/services/secauthnbproto3'
        , true);

        static::$is_initialized = true;
    }
}
