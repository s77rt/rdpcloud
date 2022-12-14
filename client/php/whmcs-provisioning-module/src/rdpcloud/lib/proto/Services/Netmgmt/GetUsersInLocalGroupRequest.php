<?php
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: services/netmgmt/local_group.proto

namespace Services\Netmgmt;

use Google\Protobuf\Internal\GPBType;
use Google\Protobuf\Internal\RepeatedField;
use Google\Protobuf\Internal\GPBUtil;

/**
 * Generated from protobuf message <code>services.netmgmt.GetUsersInLocalGroupRequest</code>
 */
class GetUsersInLocalGroupRequest extends \Google\Protobuf\Internal\Message
{
    /**
     * Generated from protobuf field <code>.models.netmgmt.LocalGroup_1 local_group = 1;</code>
     */
    protected $local_group = null;

    /**
     * Constructor.
     *
     * @param array $data {
     *     Optional. Data for populating the Message object.
     *
     *     @type \Models\Netmgmt\LocalGroup_1 $local_group
     * }
     */
    public function __construct($data = NULL) {
        \GPBMetadata\Services\Netmgmt\LocalGroup::initOnce();
        parent::__construct($data);
    }

    /**
     * Generated from protobuf field <code>.models.netmgmt.LocalGroup_1 local_group = 1;</code>
     * @return \Models\Netmgmt\LocalGroup_1|null
     */
    public function getLocalGroup()
    {
        return $this->local_group;
    }

    public function hasLocalGroup()
    {
        return isset($this->local_group);
    }

    public function clearLocalGroup()
    {
        unset($this->local_group);
    }

    /**
     * Generated from protobuf field <code>.models.netmgmt.LocalGroup_1 local_group = 1;</code>
     * @param \Models\Netmgmt\LocalGroup_1 $var
     * @return $this
     */
    public function setLocalGroup($var)
    {
        GPBUtil::checkMessage($var, \Models\Netmgmt\LocalGroup_1::class);
        $this->local_group = $var;

        return $this;
    }

}

