<?php
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: models/netmgmt/local_group.proto

namespace Models\Netmgmt;

use Google\Protobuf\Internal\GPBType;
use Google\Protobuf\Internal\RepeatedField;
use Google\Protobuf\Internal\GPBUtil;

/**
 * Generated from protobuf message <code>models.netmgmt.LocalGroup</code>
 */
class LocalGroup extends \Google\Protobuf\Internal\Message
{
    /**
     * Generated from protobuf field <code>string groupname = 1;</code>
     */
    protected $groupname = '';

    /**
     * Constructor.
     *
     * @param array $data {
     *     Optional. Data for populating the Message object.
     *
     *     @type string $groupname
     * }
     */
    public function __construct($data = NULL) {
        \GPBMetadata\Models\Netmgmt\LocalGroup::initOnce();
        parent::__construct($data);
    }

    /**
     * Generated from protobuf field <code>string groupname = 1;</code>
     * @return string
     */
    public function getGroupname()
    {
        return $this->groupname;
    }

    /**
     * Generated from protobuf field <code>string groupname = 1;</code>
     * @param string $var
     * @return $this
     */
    public function setGroupname($var)
    {
        GPBUtil::checkString($var, True);
        $this->groupname = $var;

        return $this;
    }

}

