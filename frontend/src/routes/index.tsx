import { createBrowserRouter } from "react-router-dom";
import { LoginPage, loginAction } from "../routes/login/index";
import { AccountActivation } from "./account_activation/AccountActivation";
import * as backoffice from "./backoffice/index";
import { Logout } from "./logout";
import * as studentdash from "./student_dashboard/index";

export const router = createBrowserRouter([
	{
		path: "/",
		element: <LoginPage />,
		action: loginAction,
	},
	{
		path: "/login",
		element: <LoginPage />,
		action: loginAction,
	},
	{
		path: "/logout",
		element: <Logout />,
	},
	{
		path: "/account-activation",
		element: <AccountActivation />,
	},
	...backoffice.router,
	...studentdash.router,
]);
