<?php
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: services/fileio/disk_management.proto

namespace Services\Fileio;

use Google\Protobuf\Internal\GPBType;
use Google\Protobuf\Internal\RepeatedField;
use Google\Protobuf\Internal\GPBUtil;

/**
 * Generated from protobuf message <code>services.fileio.SetQuotaStateRequest</code>
 */
class SetQuotaStateRequest extends \Google\Protobuf\Internal\Message
{
    /**
     * Generated from protobuf field <code>string volume_path = 1;</code>
     */
    protected $volume_path = '';
    /**
     * Generated from protobuf field <code>uint32 quota_state = 2;</code>
     */
    protected $quota_state = 0;

    /**
     * Constructor.
     *
     * @param array $data {
     *     Optional. Data for populating the Message object.
     *
     *     @type string $volume_path
     *     @type int $quota_state
     * }
     */
    public function __construct($data = NULL) {
        \GPBMetadata\Services\Fileio\DiskManagement::initOnce();
        parent::__construct($data);
    }

    /**
     * Generated from protobuf field <code>string volume_path = 1;</code>
     * @return string
     */
    public function getVolumePath()
    {
        return $this->volume_path;
    }

    /**
     * Generated from protobuf field <code>string volume_path = 1;</code>
     * @param string $var
     * @return $this
     */
    public function setVolumePath($var)
    {
        GPBUtil::checkString($var, True);
        $this->volume_path = $var;

        return $this;
    }

    /**
     * Generated from protobuf field <code>uint32 quota_state = 2;</code>
     * @return int
     */
    public function getQuotaState()
    {
        return $this->quota_state;
    }

    /**
     * Generated from protobuf field <code>uint32 quota_state = 2;</code>
     * @param int $var
     * @return $this
     */
    public function setQuotaState($var)
    {
        GPBUtil::checkUint32($var);
        $this->quota_state = $var;

        return $this;
    }

}

