import { MantineProvider } from '@mantine/core';
import "@mantine/core/styles.css";
import React from "react";
import ReactDOM from "react-dom/client";
import { RouterProvider, createBrowserRouter } from "react-router-dom";
import App from "./App.tsx";
import Login, { loginAction } from "./routes/Login.tsx";
import UserManagementLayout, {
	CreateManagedUser,
	ListManagedUsers,
	actionCreateManagedUser,
	loaderUserManagement,
} from "./routes/UserManagement.tsx";

const router = createBrowserRouter([
	{
		path: "/",
		element: <App />,
	},
	{
		path: "/login",
		element: <Login />,
		action: loginAction,
	},
	{
		path: "/user-management",
		element: <UserManagementLayout />,
		children: [
			{
				path: "",
				index: true,
				element: <ListManagedUsers />,
				loader: loaderUserManagement,
			},
			{
				path: "create",
				element: <CreateManagedUser />,
				action: actionCreateManagedUser,
			},
		],
	},
]);

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
	<React.StrictMode>
		<MantineProvider>
			<RouterProvider router={router} />
		</MantineProvider>
	</React.StrictMode>,
);
