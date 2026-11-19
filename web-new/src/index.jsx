/* @refresh reload */
import MainPage from "./routes/MainPage.jsx";
import ProfilePage from "./routes/ProfilePage.jsx";
import ChallengePage from "./routes/ChallengePage.jsx";
import AdminPage from "./routes/AdminPage.jsx";
import LoginPage from "./routes/LoginPage.jsx";
import LogoutPage from "./routes/Logout.jsx";
import ForgotPage from "./routes/ForgotPage.jsx";
import ResetPage from "./routes/ResetPage.jsx";

import MainLayout from "./layouts/MainLayout";
import BasicLayout from "./layouts/BasicLayout";

import ProtectedRoute from "./components/ProtectedRoute";
import AnonymousRoute from "./components/AnonymousRoute";
import { AuthProvider } from "./context/AuthContext";

import { render } from 'solid-js/web'
import { Router, Route } from "@solidjs/router";
import './index.css'

const root = document.getElementById('root')

render(() =>
		<Router>
				<Route component={AuthProvider}>
						<Route component={AnonymousRoute}>
								<Route component={BasicLayout}>
										<Route
												path="/login"
												component={LoginPage}
										/>
										<Route
												path="/forgot"
												component={ForgotPage}
										/>
										<Route
												path="/reset-password"
												component={ResetPage}
										/>
								</Route>
						</Route>
						<Route component={ProtectedRoute}>
                <Route component={MainLayout}>
                    <Route
                        path="/"
                        component={MainPage}
                    />
                    <Route
                        path="/profile"
                        component={ProfilePage}
                    />
                    <Route
                        path="/challenges"
                        component={ChallengePage}
                    />
                    <Route
                        path="/admin"
                        component={AdminPage}
                    />
										<Route
                        path="/logout"
                        component={LogoutPage}
                    />
                </Route>
						</Route>
				</Route>
		</Router>
		, root
)
