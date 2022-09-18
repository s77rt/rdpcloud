import React from 'react';
import './App.css';
import { ConfigProvider, Layout, Form, Input, Button, message, notification, Space, Typography, Result, Menu, Alert, Spin } from 'antd';
import { SmileOutlined, UserOutlined, LockOutlined, LoadingOutlined } from '@ant-design/icons';
import { parseJWT, isExpiredJWT } from "./jwt";

ConfigProvider.config({
	theme: {
		primaryColor: '#00c292',
	},
});

const { Header, Content, Footer, Sider } = Layout;
const { Title, Link } = Typography;

class App extends React.Component {
	constructor() {
		super();

		this.onClickBrand = this.onClickBrand.bind(this);
		this.onFinishLogin = this.onFinishLogin.bind(this);
		this.onClickLogout = this.onClickLogout.bind(this);
		this.onCollapseSider = this.onCollapseSider.bind(this);
		this.onClickInvoke = this.onClickInvoke.bind(this);
		this.onSelectMenu = this.onSelectMenu.bind(this);

		this.state = {
			appName: "RDPCloud",

			services: {},

			showLogin: true,
			showRequestBuilder: false,

			selectedKeysMenu: [],

			sider_isCollapsed: false,
			invoke_isLoading: false,
			login_isLoading: false,
			selecting_service: false,
			selecting_method: false,

			response_error: null,
			response_data: null,

			login_error: null,

			token: null,
			user: undefined
		}
	}

	onClickBrand() {
		this.setState({
			showRequestBuilder: false,
			selectedKeysMenu: []
		});
	}

	onFinishLogin(values) {
		this.Login(values);
	}

	onClickLogout() {
		this.Logout();
	}

	onCollapseSider(collapsed) {
		this.setState({
			sider_isCollapsed: collapsed,
		})
	}

	onClickInvoke() {
		this.Invoke();
	}

	onSelectMenu(item) {
		if (item.keyPath.length !== 2)
			return;
		this.SelectService(item.keyPath[1]);
		this.SelectMethod(item.keyPath[0]);
		this.setState({
			selectedKeysMenu: item.selectedKeys,
			showRequestBuilder: true,
			response_error: null,
			response_data: null
		});
	}

	SetAppName(name) {
		this.setState({
			appName: name
		});
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

	SetSelectingService(selecting) {
		this.setState({
			selecting_service: selecting,
		})
	}

	SetSelectingMethod(selecting) {
		this.setState({
			selecting_method: selecting,
		})
	}

	SelectService(key) {
		this.SetSelectingService(true);
		document.dispatchEvent(new CustomEvent("App-SelectService", { detail: key }));
	}

	SelectMethod(key) {
		this.SetSelectingMethod(true);
		document.dispatchEvent(new CustomEvent("App-SelectMethod", { detail: key }));
	}

	SetInvokeLoading(loading) {
		this.setState({
			invoke_isLoading: loading,
		})
	}

	InvokeResponseCallback(data, error, error_msg) {
		if (error || error_msg) {
			this.setState({
				response_error: {
					error: error || "",
					error_msg: error_msg || "",
				}
			})
			return;
		}

		try {
			data = JSON.parse(data);
		} catch (e) {
			this.ThrowErrorMessage("Invalid data: " + e.message);
			return;
		}

		this.setState({
			response_data: {
				data: data
			}
		});
	}

	Invoke() {
		this.SetInvokeLoading(true);
		this.setState({
			response_error: null,
			response_data: null
		});

		var jwt;

		try {
			jwt = parseJWT(this.state.token);
		} catch (e) {
			this.ThrowErrorMessage("Invalid token: " + e.message);
			this.SetInvokeLoading(false);
			return;
		}

		if (isExpiredJWT(jwt)) {
			this.Logout();
			this.PushInfoNotification("Your token has expired", "To continue using the service please login.");
			this.SetInvokeLoading(false);
			return;
		}

		document.dispatchEvent(new CustomEvent("App-Invoke", { detail: this.state.token }));
	}

	SetLoginLoading(loading) {
		this.setState({
			login_isLoading: loading,
		})
	}

	LoginResponseCallback(data, error, error_msg) {
		if (error || error_msg) {
			this.setState({
				login_error: {
					error: error || "",
					error_msg: error_msg || "",
				}
			})
			return;
		}

		try {
			data = JSON.parse(data);
		} catch (e) {
			this.ThrowErrorMessage("Invalid data: " + e.message);
			return;
		}

		var jwt;

		try {
			jwt = parseJWT(data.token);
		} catch (e) {
			this.ThrowErrorMessage("Invalid token: " + e.message);
			return;
		}

		this.setState({
			showLogin: false,
			token: data.token,
			user: {preferred_username: jwt.preferred_username}
		});
	}

	Login(credentials) {
		this.SetLoginLoading(true);
		this.setState({
			login_error: null
		});
		document.dispatchEvent(new CustomEvent("App-Login", { detail: credentials }));
	}

	Logout() {
		this.setState({
			showLogin: true,
			showRequestBuilder: false,
			selectedKeysMenu: [],
			token: null,
			user: undefined
		});
	}

	PushInfoNotification(message, description) {
		notification.info({
			message: message,
			description: description,
			placement: "top"
		});
	}

	ThrowErrorMessage(msg) {
		message.error(msg);
	}

	componentDidMount() {
		if (typeof window.InitApp === "function")
			window.InitApp(this); // async function
	}

	render() {
		return (
			<Layout style={{ minHeight: "100vh" }}>
				<Header className="header" style={{ display: "flex", justifyContent: "space-between", position: "fixed", zIndex: 1, width: "100%", whiteSpace: "nowrap", overflow: "auto" }}>
					<strong className="brand" onClick={this.onClickBrand}>{this.state.appName + " Web Client"}</strong>
					{this.state.user &&
						<Space size="middle">
							<span>Logged in as <strong>{this.state.user.preferred_username}</strong></span>
							<Button onClick={this.onClickLogout} type="primary" className="logout-btn">Logout</Button>
						</Space>
					}
				</Header>
				<Layout style={{ marginTop: "64px" }}>
					<Sider
						width={200}
						breakpoint="sm"
						collapsedWidth={this.state.showLogin === true ? 0 : 50}
						collapsible
						collapsed={this.state.showLogin === true ? true : this.state.sider_isCollapsed}
						onCollapse={this.onCollapseSider}
						style={{
							overflow: "auto",
							height: "100vh",
							position: "fixed",
							left: 0,
							top: "64px",
							bottom: 0,
							cursor: (this.state.selecting_service || this.state.selecting_method) ? "not-allowed" : "auto"
						}}
					>
						<Menu
							selectable={!(this.state.selecting_service || this.state.selecting_method)}
							onSelect={this.onSelectMenu}
							selectedKeys={this.state.selectedKeysMenu}
							mode="vertical"
							style={{
								height: "100%",
								borderRight: 0,
								pointerEvents: (this.state.selecting_service || this.state.selecting_method) ? "none" : "auto"
							}}
							items={Object.values(this.state.services)}
						/>
					</Sider>
					<Layout
						style={{
							marginLeft: this.state.showLogin === true ? 0 : this.state.sider_isCollapsed === true ? "50px" : "200px",
							transition: "all .2s",
							padding: "0 24px"
						}}
					>
						<Content style={{
							padding: 24,
							margin: 0,
							minHeight: 280
						}}>
							{this.state.showLogin === true &&
								<>
									<Result
										status="info"
										icon=<LockOutlined />
										title="Login"
										subTitle="Access is restricted to authorized users only"
										extra={this.state.login_error ?
											<Alert
												message={this.state.login_error.error + ": " + this.state.login_error.error_msg}
												type="error"
												showIcon
												style={{ maxWidth: "320px", margin: "auto" }}
											/>
										: null}
									/>
									<div className="login-form-container">
										<Form
											name="login"
											className="login-form"
											autoComplete="off"
											size="large"
											onFinish={this.onFinishLogin}
										>
											<Form.Item
												name="username"
												rules={[{ required: true, message: "Username is required" }]}
											>
												<Input
													prefix=<UserOutlined className="input-icon" />
													placeholder="Username"
													disabled={this.state.login_isLoading}
												/>
											</Form.Item>
											<Form.Item
												name="password"
												rules={[{ required: true, message: "Password is required" }]}
											>
												<Input
													prefix=<LockOutlined className="input-icon" />
													type="password"
													placeholder="Password"
													disabled={this.state.login_isLoading}
												/>
											</Form.Item>
											<Form.Item>
												<Button loading={this.state.login_isLoading} type="primary" htmlType="submit" className="login-form-button">
													Log in
												</Button>
											</Form.Item>
										</Form>
									</div>
								</>
							}
							<Space direction="vertical" size="middle" style={{ width: "100%", display: this.state.showRequestBuilder === true ? "flex" : "none"}}>
								<div>
									<Title level={4}>Request</Title>
									<div className="content-box">
										{(this.state.selecting_service || this.state.selecting_method) &&
											<Spin indicator={<LoadingOutlined/>} />
										}
										<div
											id="App-Request-Form-Container"
											style={{ display: (this.state.selecting_service || this.state.selecting_method) ? "none" : "block" }}
										></div>
										<Button
											type="primary"
											style={{ display: (this.state.selecting_service || this.state.selecting_method) ? "none" : "inline-block" }}
											loading={this.state.invoke_isLoading}
											onClick={this.onClickInvoke}>
											Invoke
										</Button>
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
							{this.state.showRequestBuilder === false && this.state.showLogin === false &&
								<Result
									status="info"
									icon={<SmileOutlined />}
									title="Ready when you are"
									subTitle="To get started choose from the side menu"
								/>
							}
						</Content>
						<Footer style={{ textAlign: "center" }}>
							<Space direction="vertical" size="small">
								<Link href="https://rdpcloud.io">Powered by RDPCloud.io</Link>
								{"Copyright © " + new Date().getFullYear() + " " + this.state.appName + ". All Rights Reserved"}
							</Space>
						</Footer>
					</Layout>
				</Layout>
			</Layout>
		);
	}
}

export default App;
