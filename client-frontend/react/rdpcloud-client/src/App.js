import React from 'react';
import './App.css';
import { ConfigProvider, Layout, Button, message, Space, Typography, Result, Menu } from 'antd';
import { SmileOutlined } from '@ant-design/icons';

ConfigProvider.config({
	theme: {
		primaryColor: '#00c292',
	},
});

const { Header, Content, Footer, Sider } = Layout;
const { Title } = Typography;

const invokeEvent = new Event('App-Invoke');

class App extends React.Component {
	constructor() {
		super();

		this.onCollapseSider = this.onCollapseSider.bind(this);
		this.onClickInvoke = this.onClickInvoke.bind(this);
		this.onSelectMenu = this.onSelectMenu.bind(this);

		this.state = {
			services: {},

			showContent: false,

			sider_isCollapsed: false,
			invoke_isLoading: false,

			response_error: null,
			response_data: null,
		}
	}

	onCollapseSider(collapsed) {
		this.setState({
			sider_isCollapsed: collapsed,
		})
	}

	onClickInvoke() {
		this.SetInvokeLoading(true);
		this.setState({
			response_error: null,
			response_data: null
		}, () => {
			document.dispatchEvent(invokeEvent);
		});
	}

	onSelectMenu(item) {
		if (item.keyPath.length !== 2)
			return;
		this.SelectService(item.keyPath[1]);
		this.SelectMethod(item.keyPath[0]);
	}

	AddService(service) {
		this.setState(previousState => ({
			services: {
				...previousState.services,
				[service]: {
					label: service.split(".").pop(),
					key: service,
					methods: {},
					children: []
				}
			}
		}));
	}

	AddServiceMethod(service, method) {
		this.setState({}, () => {
			if (!(service in this.state.services))
				return
			this.setState(previousState => ({
				services: {
					...previousState.services,
					[service]: {
						...previousState.services[service],
						methods: {
							...previousState.services[service].methods,
							[method]: {
								label: method.replace(/([A-Z])/g, ' $1').trim(),
								key: method
							}
						},
						children: []
					}
				}
			}), () => {
				this.setState(previousState => ({
					services: {
						...previousState.services,
						[service]: {
							...previousState.services[service],
							children: Object.values(previousState.services[service].methods)
						}
					}
				}));
			});
		});
	}

	SetInvokeLoading(loading) {
		this.setState({
			invoke_isLoading: loading,
		})
	}

	SelectService(key) {
		document.dispatchEvent(new CustomEvent("App-SelectService", { detail: key }));
	}

	SelectMethod(key) {
		document.dispatchEvent(new CustomEvent("App-SelectMethod", { detail: key }));
		this.setState({
			showContent: true,
			response_error: null,
			response_data: null
		})
	}

	SetResponse(data, error, error_msg) {
		if (error || error_msg) {
			this.setState({
				response_error: {
					error: error || "",
					error_msg: error_msg || "",
				},
				response_data: null
			})
			return
		}
		try {
			data = JSON.parse(data);
		} catch (e) {
			this.ThrowErrorMessage(e.message);
			return;
		}
		this.setState({
			response_error: null,
			response_data: {
				data: data
			}
		});
	}

	ThrowErrorMessage(msg) {
		message.error(msg);
	}

	componentDidMount() {
		if (typeof window.InitApp === "function")
			window.InitApp(this);
	}

	render() {
		return (
			<Layout style={{ minHeight: "100vh" }}>
				<Header className="header" style={{ position: "fixed", zIndex: 1, width: "100%" }}>
					<span>RDPCloud Web Client</span>
				</Header>
				<Layout style={{ marginTop: "64px" }}>
					<Sider
						width={200}
						breakpoint="sm"
						collapsedWidth={50}
						collapsible
						collapsed={this.state.sider_isCollapsed}
						onCollapse={this.onCollapseSider}
						style={{
							overflow: "auto",
							height: "100vh",
							position: "fixed",
							left: 0,
							top: "64px",
							bottom: 0,
					  	}}
					>
						<Menu
							onSelect={this.onSelectMenu}
							mode="vertical"
							style={{ height: "100%", borderRight: 0 }}
							items={Object.values(this.state.services)}
						/>
					</Sider>
					<Layout
						style={{
							marginLeft: this.state.sider_isCollapsed ? "50px" : "200px",
							transition: "all .2s",
							padding: "0 24px"
						}}
					>
						<Content style={{
							padding: 24,
							margin: 0,
							minHeight: 280
						}}>
							<Space direction="vertical" size="middle" style={{ display: this.state.showContent === true ? "flex" : "none"}}>
								<div>
									<Title level={4}>Request</Title>
									<div className="content-box">
										<div id="App-Request-Form-Container"></div>
										<Button type="primary" loading={this.state.invoke_isLoading} onClick={this.onClickInvoke}>Invoke</Button>
									</div>
								</div>
								<div>
									<Title level={4}>Response</Title>
									<div className="content-box">
										{this.state.response_error ? (
											<Result
												status="error"
												title={this.state.response_error.error}
												subTitle={this.state.response_error.error_msg}
											/>
										) : (this.state.response_data ? (
											(Object.keys(this.state.response_data.data).length === 0 ? (
												<Result
													status="success"
													title="Operation completed successfully"
													subTitle="Your operation has been executed"
												/>
											) : (
												<textarea>{JSON.stringify(this.state.response_data.data, null, 2)}</textarea>
											))
										) : null
										)}
									</div>
								</div>
							</Space>
							{this.state.showContent === false &&
								<Result
									status="info"
									icon={<SmileOutlined />}
									title="Ready when you are"
									subTitle="To get started choose from the side menu"
								/>
							}
						</Content>
						<Footer style={{ textAlign: "center" }}>
							RDPCloud
						</Footer>
					</Layout>
				</Layout>
			</Layout>
		);
	}
}

export default App;
