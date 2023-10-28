import React from "react";
import ReactDOM from "react-dom/client";
import { RouterProvider, createBrowserRouter } from "react-router-dom";
import App from "./App.tsx";
import Login, { loginAction } from "./routes/Login.tsx";
import UserManagement, { CreateManagedUser, ListManagedUsers, actionCreateManagedUser, loaderUserManagement } from "./routes/UserManagement.tsx";

const router = createBrowserRouter([
	{ 
		path: "/", 
		element: <App />
	},
	{
		path: "/login",
		element: <Login />,
		action: loginAction,
	},
	{
		path: "/user-management",
		element: <UserManagement />,
		children: [
			{
				path: "",
				element: <ListManagedUsers />,
				loader: loaderUserManagement,
			},
		]
	},
	{
		path: "/user-management/create",
		element: <CreateManagedUser />,
		action: actionCreateManagedUser,
	},
]);

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
	<React.StrictMode>
		<RouterProvider router={router} />
	</React.StrictMode>,
);
