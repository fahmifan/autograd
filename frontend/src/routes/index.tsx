import { createBrowserRouter } from "react-router-dom";
import { LoginPage, loginAction } from "../routes/login/index";
import * as backoffice from "./backoffice/index";
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
	...backoffice.router,
	...studentdash.router,
]);