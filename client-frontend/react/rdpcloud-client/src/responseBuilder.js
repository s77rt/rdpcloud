import { Result, Empty, Alert, Typography, Table, Statistic, Row, Col } from 'antd';

const { Paragraph } = Typography;

const isArray = (o) => { return o != null && typeof(o) === "object" && Array.isArray(o) }
const isObject = (o) => { return o != null && typeof(o) === "object" && !Array.isArray(o) }
const isString = (o) => { return typeof(o) === "string" }
const isNumber = (o) => { return typeof(o) === "number" }

export const buildResponseFromResponseData = (response_data) => {
	if (isObject(response_data)) {
		const keys = Object.keys(response_data);
		if (keys.length === 0) {
			return (
				<Result
					status="success"
					title="Operation completed successfully"
					subTitle="Your operation has been executed"
				/>
			);
		} else if (keys.length === 1) {
			const records = response_data[keys[0]];
			if (isObject(records)) {
				return (
					<Row gutter={16}>
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
					<Row gutter={16}>
						<Col span={24}>
							<Statistic title={keys[0].replace(/_/g, " ").toUpperCase().trim()} value={records} />
						</Col>
					</Row>
				)
			}
		}
	}

	// fallback
	if (isObject(response_data)) {
		return (
			<Paragraph code>
				<pre>
					{JSON.stringify(response_data, null, 2)}
				</pre>
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
