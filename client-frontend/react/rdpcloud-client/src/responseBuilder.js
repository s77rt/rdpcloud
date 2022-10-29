import { Result, Empty, Alert, Typography, Table, Statistic, Row, Col, Badge, Tree } from 'antd';
import { CheckCircleOutlined, CloseCircleOutlined } from '@ant-design/icons';
import { red, green } from '@ant-design/colors';
import prettyBytes from 'pretty-bytes';
import * as dayjs from 'dayjs'
import * as relativeTime from 'dayjs/plugin/relativeTime' // import plugin
dayjs.extend(relativeTime);

const { Paragraph } = Typography;

const isArray = (o) => { return o != null && typeof(o) === "object" && Array.isArray(o) }
const isObject = (o) => { return o != null && typeof(o) === "object" && !Array.isArray(o) }
const isString = (o) => { return typeof(o) === "string" }
const isNumber = (o) => { return typeof(o) === "number" }

export const buildResponseFromResponseData = (response_data) => {
	const service = response_data.service;
	const method = response_data.method;
	const data = response_data.data;

	if (`/${service}/${method}` in builderMap) {
		try {
			return builderMap[`/${service}/${method}`](data);
		} catch {
			return genericResponseBuilder(data);
		}
	}

	return genericResponseBuilder(data);
}

function genericResponseBuilder(data) {
	if (isObject(data)) {
		const keys = Object.keys(data);
		if (keys.length === 0) {
			return (
				<Result
					status="success"
					title="Operation completed successfully"
					subTitle="Your operation has been executed"
				/>
			);
		} else if (keys.length === 1) {
			const records = data[keys[0]];
			if (isObject(records)) {
				return (
					<Row gutter={[16, 14]}>
						{Object.entries(records).map(([k, v], i) => (
							<Col span={12}>
								<Statistic key={i} title={k.replace(/_/g, " ").toUpperCase().trim()} value={v} />
							</Col>
						))}
					</Row>
				)
			} else if (isArray(records)) {
				if (records.length === 0) {
					return (
						<Empty />
					)
				}
				var first_record = records[0];
				if (isObject(first_record)) {
					const columns = Object.keys(first_record).map((v) => (
						{
							title: v.replace(/_/g, " ").toUpperCase().trim(),
							dataIndex: v,
							key: v
						}
					));
					const dataSource = records.map((v, i) => ({...v, key: i}));
					return (
						<Table dataSource={dataSource} columns={columns} />
					);
				}
			} else if (isString(records) || isNumber(records)) {
				return (
					<Row gutter={[16, 14]}>
						<Col span={24}>
							<Statistic title={keys[0].replace(/_/g, " ").toUpperCase().trim()} value={records} />
						</Col>
					</Row>
				)
			}
		}
	}

	// fallback
	if (isObject(data)) {
		return (
			<Paragraph className="code-block-wrapper" code>
				{JSON.stringify(data, null, 2)}
			</Paragraph>
		);
	}
	return (
		<Alert
			message="Unknown response data structure"
			description="Unable to process the response data"
			type="warning"
		/>
	);
}

const builderMap = {
	"/services.sysinfo.Sysinfo/GetUptime": (data) => {
		const keys = Object.keys(data);
		const records = data[keys[0]];

		const record = dayjs(0).to(dayjs(Number(records)), true);
		return (
			<Row gutter={[16, 14]}>
				<Col span={24}>
					<Statistic title={keys[0].replace(/_/g, " ").toUpperCase().trim()} value={record} />
				</Col>
			</Row>
		)
	},
	"/services.netmgmt.Netmgmt/GetUsers": (data) => {
		const UF_ACCOUNTDISABLE = 0x0002;
		const USER_PRIV_GUEST = 0
		const USER_PRIV_USER = 1
		const USER_PRIV_ADMIN = 2

		const userPrivMeaning = {
			[USER_PRIV_GUEST]: "Guest",
			[USER_PRIV_USER]: "User",
			[USER_PRIV_ADMIN]: "Administrator"
		}

		const keys = Object.keys(data);
		const records = data[keys[0]];

		if (records.length === 0) {
			return (
				<Empty />
			)
		}
		var first_record = records[0];

		var columns = [];
		for (let v of [...Object.keys(first_record), "active"]) {
			if (["password", "flags"].indexOf(v) >= 0)
				continue;
			let column = {
				title: v.replace(/_/g, " ").toUpperCase().trim(),
				dataIndex: v,
				align: v === "active" ? "center" : "left",
				sorter: (() => {
					let sorterFunc;
					switch (v) {
						case "username":
						case "privilege":
							sorterFunc = (a, b) => { return a[v].localeCompare(b[v], 'en', { numeric: true }) };
							break;
						case "active":
							sorterFunc = (a, b) => { return a[v].props["data-value"] - b[v].props["data-value"] };
							break;
						default:
							sorterFunc = undefined;
					}
					return sorterFunc;
				})(),
				key: v
			}
			columns.push(column);
		}

		var dataSource = [];
		records.forEach(function (v, i) {
			let record = {
				...v,
				"privilege": userPrivMeaning[v.privilege],
				"active": (v.flags & UF_ACCOUNTDISABLE) === 0 ? <CheckCircleOutlined data-value={true} style={{color: green[6]}} /> : <CloseCircleOutlined data-value={false} style={{color: red[6]}} />,
				key: i
			}
			delete record.password;
			delete record.flags;
			dataSource.push(record);
		});

		return (
			<Table dataSource={dataSource} columns={columns} />
		);
	},
	"/services.netmgmt.Netmgmt/GetUser": (data) => {
		const UF_ACCOUNTDISABLE = 0x0002;
		const USER_PRIV_MASK = 0x3;
		const USER_PRIV_GUEST = 0
		const USER_PRIV_USER = 1
		const USER_PRIV_ADMIN = 2

		const userPrivMeaning = {
			[USER_PRIV_GUEST]: "Guest",
			[USER_PRIV_USER]: "User",
			[USER_PRIV_ADMIN]: "Administrator"
		}

		const keys = Object.keys(data);
		const records = data[keys[0]];

		const v = records;
		let record = {
			...v,
			"privilege": userPrivMeaning[v.privilege],
			"active": (v.flags & UF_ACCOUNTDISABLE) === 0 ? "Yes" : "No",
		}
		delete record.password;
		delete record.flags;

		return (
			<Row gutter={[16, 14]}>
				{Object.entries(record).map(([k, v], i) => (
					<Col span={12}>
						<Statistic key={i} title={k.replace(/_/g, " ").toUpperCase().trim()} value={v} />
					</Col>
				))}
			</Row>
		);
	},
	"/services.netmgmt.Netmgmt/GetMyUser": function (data) {
		return this["/services.netmgmt.Netmgmt/GetUser"](data);
	},
	"/services.fileio.Fileio/GetQuotaState": (data) => {
		const DISKQUOTA_STATE_DISABLED = 0x00000000;
		const DISKQUOTA_STATE_TRACK = 0x00000001;
		const DISKQUOTA_STATE_ENFORCE = 0x00000002;
		const DISKQUOTA_STATE_MASK = 0x00000003;
		const DISKQUOTA_FILESTATE_INCOMPLETE = 0x00000100;
		const DISKQUOTA_FILESTATE_REBUILDING = 0x00000200;
		const DISKQUOTA_FILESTATE_MASK = 0x00000300;

		const quotaStateMeaning = {
			[DISKQUOTA_STATE_DISABLED]: "Quotas are not enabled on the volume",
			[DISKQUOTA_STATE_TRACK]: "Quotas are enabled but the limit value is not being enforced. Users may exceed their quota limit",
			[DISKQUOTA_STATE_ENFORCE]: "Quotas are enabled and the limit value is enforced. Users cannot exceed their quota limit",
			[DISKQUOTA_FILESTATE_INCOMPLETE]: "The volume's quota information is out of date. Quotas are probably disabled",
			[DISKQUOTA_FILESTATE_REBUILDING]: "The volume is rebuilding its quota information"
		}

		const keys = Object.keys(data);
		const records = data[keys[0]];

		const record = records;
		const record_meaning = quotaStateMeaning[record];
		return (
			<Row gutter={[16, 14]}>
				<Col span={24}>
					<Statistic title={keys[0].replace(/_/g, " ").toUpperCase().trim()} value={record} />
				</Col>
				<Col span={24}>
					<Statistic title={(keys[0]+"_meaning").replace(/_/g, " ").toUpperCase().trim()} value={record_meaning} />
				</Col>
			</Row>
		)
	},
	"/services.fileio.Fileio/GetDefaultQuota": (data) => {
		const keys = Object.keys(data);
		const records = data[keys[0]];

		const v = records;
		let record = {
			...v,
			"quota_threshold": prettyBytes(Number(v.quota_threshold), {binary: true}),
			"quota_limit": prettyBytes(Number(v.quota_limit), {binary: true}),
		}

		return (
			<Row gutter={[16, 14]}>
				{Object.entries(record).map(([k, v], i) => (
					<Col span={12}>
						<Statistic key={i} title={k.replace(/_/g, " ").toUpperCase().trim()} value={v} />
					</Col>
				))}
			</Row>
		);
	},
	"/services.fileio.Fileio/GetUsersQuotaEntries": (data) => {
		const DISKQUOTA_USER_ACCOUNT_RESOLVED = 0;
		const DISKQUOTA_USER_ACCOUNT_UNAVAILABLE = 1;
		const DISKQUOTA_USER_ACCOUNT_DELETED = 2;
		const DISKQUOTA_USER_ACCOUNT_INVALID = 3;
		const DISKQUOTA_USER_ACCOUNT_UNKNOWN = 4;
		const DISKQUOTA_USER_ACCOUNT_UNRESOLVED = 5;

		const accountStatusMeaning = {
			[DISKQUOTA_USER_ACCOUNT_RESOLVED]: "Resolved",
			[DISKQUOTA_USER_ACCOUNT_UNAVAILABLE]: "Unavailable",
			[DISKQUOTA_USER_ACCOUNT_DELETED]: "Deleted",
			[DISKQUOTA_USER_ACCOUNT_INVALID]: "Invalid",
			[DISKQUOTA_USER_ACCOUNT_UNKNOWN]: "Unknown",
			[DISKQUOTA_USER_ACCOUNT_UNRESOLVED]: "Unresolved"
		}
		const accountStatusMeaning_Status = {
			[DISKQUOTA_USER_ACCOUNT_RESOLVED]: "success",
			[DISKQUOTA_USER_ACCOUNT_UNAVAILABLE]: "warning",
			[DISKQUOTA_USER_ACCOUNT_DELETED]: "error",
			[DISKQUOTA_USER_ACCOUNT_INVALID]: "error",
			[DISKQUOTA_USER_ACCOUNT_UNKNOWN]: "default",
			[DISKQUOTA_USER_ACCOUNT_UNRESOLVED]: "default"
		}

		const keys = Object.keys(data);
		const records = data[keys[0]];

		if (records.length === 0) {
			return (
				<Empty />
			)
		}
		var first_record = records[0];

		var columns = [];
		for (let v of Object.keys(first_record)) {
			let column = {
				title: v.replace(/_/g, " ").toUpperCase().trim(),
				dataIndex: v,
				sorter: (() => {
					let sorterFunc;
					switch (v) {
						case "username":
							sorterFunc = (a, b) => { return a[v].localeCompare(b[v], 'en', { numeric: true }) };
							break;
						case "account_status":
						case "quota_threshold":
						case "quota_limit":
						case "quota_used":
							sorterFunc = (a, b) => { return a[v].props["data-value"] - b[v].props["data-value"] };
							break;
						default:
							sorterFunc = undefined;
					}
					return sorterFunc;
				})(),
				key: v
			}
			columns.push(column);
		}

		var dataSource = [];
		records.forEach(function (v, i) {
			let record = {
				...v,
				"account_status": <Badge data-value={v.account_status} status={accountStatusMeaning_Status[v.account_status]} text={accountStatusMeaning[v.account_status]} />,
				"quota_threshold": <span data-value={Number(v.quota_threshold)}>{prettyBytes(Number(v.quota_threshold), {binary: true})}</span>,
				"quota_limit": <span data-value={Number(v.quota_limit)}>{prettyBytes(Number(v.quota_limit), {binary: true})}</span>,
				"quota_used": <span data-value={Number(v.quota_used)}>{prettyBytes(Number(v.quota_used), {binary: true})}</span>,
				key: i
			}
			dataSource.push(record);
		});

		return (
			<Table dataSource={dataSource} columns={columns} />
		);
	},
	"/services.fileio.Fileio/GetUserQuotaEntry": (data) => {
		const DISKQUOTA_USER_ACCOUNT_RESOLVED = 0;
		const DISKQUOTA_USER_ACCOUNT_UNAVAILABLE = 1;
		const DISKQUOTA_USER_ACCOUNT_DELETED = 2;
		const DISKQUOTA_USER_ACCOUNT_INVALID = 3;
		const DISKQUOTA_USER_ACCOUNT_UNKNOWN = 4;
		const DISKQUOTA_USER_ACCOUNT_UNRESOLVED = 5;

		const accountStatusMeaning = {
			[DISKQUOTA_USER_ACCOUNT_RESOLVED]: "Resolved",
			[DISKQUOTA_USER_ACCOUNT_UNAVAILABLE]: "Unavailable",
			[DISKQUOTA_USER_ACCOUNT_DELETED]: "Deleted",
			[DISKQUOTA_USER_ACCOUNT_INVALID]: "Invalid",
			[DISKQUOTA_USER_ACCOUNT_UNKNOWN]: "Unknown",
			[DISKQUOTA_USER_ACCOUNT_UNRESOLVED]: "Unresolved"
		}

		const keys = Object.keys(data);
		const records = data[keys[0]];

		const v = records;
		let record = {
			...v,
			"account_status": accountStatusMeaning[v.account_status],
			"quota_threshold": prettyBytes(Number(v.quota_threshold), {binary: true}),
			"quota_limit": prettyBytes(Number(v.quota_limit), {binary: true}),
			"quota_used": prettyBytes(Number(v.quota_used), {binary: true}),
		}

		return (
			<Row gutter={[16, 14]}>
				{Object.entries(record).map(([k, v], i) => (
					<Col span={12}>
						<Statistic key={i} title={k.replace(/_/g, " ").toUpperCase().trim()} value={v} />
					</Col>
				))}
			</Row>
		);
	},
	"/services.fileio.Fileio/GetMyUserQuotaEntry": function (data) {
		return this["/services.fileio.Fileio/GetUserQuotaEntry"](data);
	},
	"/services.fileio.Fileio/GetVolumes": (data) => {
		const keys = Object.keys(data);
		const records = data[keys[0]];

		const treeData = records.map(record => ({
			title: record.guid,
			key: record.guid,
			children: record.paths.map((path, i) => ({
				title: path,
				key: `${record.guid}@${i}`
			}))
		}))

		return (
			<Tree treeData={treeData} />
		);
	},
	"/services.msi.Msi/GetProducts": (data) => {
		const keys = Object.keys(data);
		const records = data[keys[0]];

		if (records.length === 0) {
			return (
				<Empty />
			)
		}
		var first_record = records[0];

		var columns = [];
		for (let v of Object.keys(first_record)) {
			if (["guid"].indexOf(v) >= 0)
				continue;
			let column = {
				title: v.replace(/_/g, " ").toUpperCase().trim(),
				dataIndex: v,
				sorter: (() => {
					let sorterFunc;
					switch (v) {
						case "name":
						case "version":
						case "publisher":
						case "install_date":
							sorterFunc = (a, b) => { return a[v].localeCompare(b[v], 'en', { numeric: true }) };
							break;
						default:
							sorterFunc = undefined;
					}
					return sorterFunc;
				})(),
				key: v
			}
			columns.push(column);
		}

		var dataSource = [];
		records.forEach(function (v, i) {
			let record = {
				...v,
				key: i
			}
			delete record.guid;
			dataSource.push(record);
		});

		return (
			<Table dataSource={dataSource} columns={columns} />
		);
	}
};
