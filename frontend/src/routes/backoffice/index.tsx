import { AppShell, Container, NavLink, Text } from "@mantine/core";
import { RouteObject } from "react-router-dom";
import { Link, Outlet, useLocation } from "react-router-dom";
import {
	CreateAssignment,
	DetailAssignment,
	ListAssignments,
	actionCreateAssignemnt,
	loadEditAssignment,
	loaderListAssignments,
} from "./Assignment";
import {
	CreateManagedUser,
	ListManagedUsers,
	actionCreateManagedUser,
	loaderUserManagement,
} from "./UserManagement";

export const router: RouteObject[] = [
	{
		path: "/backoffice",
		element: <DashboardLayout />,
		children: [
			{
				path: "/backoffice/user-management",
				index: true,
				element: <ListManagedUsers />,
				loader: loaderUserManagement,
			},
			{
				path: "/backoffice/user-management/create",
				element: <CreateManagedUser />,
				action: actionCreateManagedUser,
			},
			{
				path: "/backoffice/assignments",
				element: <ListAssignments />,
				loader: loaderListAssignments,
			},
			{
				path: "/backoffice/assignments/create",
				element: <CreateAssignment />,
				action: actionCreateAssignemnt,
			},
			{
				path: "/backoffice/assignments/detail",
				element: <DetailAssignment />,
				loader: loadEditAssignment,
			},
		],
	},
];

export default function DashboardLayout() {
	const location = useLocation();

	const navitems = [
		{
			label: "User Management",
			to: "/backoffice/user-management",
			children: [
				{
					label: "List Users",
					to: "/backoffice/user-management",
				},
				{
					label: "Create User",
					to: "/backoffice/user-management/create",
				},
			],
		},
		{
			label: "Assignment Submission",
			to: "/backoffice/assignments",
			children: [
				{
					label: "List Assignments",
					to: "/backoffice/assignments",
				},
				{
					label: "Create Assignment",
					to: "/backoffice/assignments/create",
				},
			],
		},
	];

	function navItemActive(path: string): boolean {
		const locpath = trimTrailingSlash(location.pathname);
		const currpath = trimTrailingSlash(path);
		return locpath === currpath;
	}

	function navItemOpened(path: string): boolean {
		const locpath = trimTrailingSlash(location.pathname);
		const currpath = trimTrailingSlash(path);
		return locpath.startsWith(currpath);
	}

	function trimTrailingSlash(path: string): string {
		if (path.trim().endsWith("/")) {
			return path.slice(0, -1);
		}
		return path;
	}

	return (
		<AppShell
			header={{ height: 60 }}
			navbar={{ width: 300, breakpoint: "sm" }}
			padding="md"
		>
			<AppShell.Header>
				<Text py="sm" px="sm" size="lg">
					Autograd Dashboard
				</Text>
			</AppShell.Header>

			<AppShell.Navbar p="md">
				{navitems.map((item) => {
					return (
						<NavLink
							key={item.to}
							label={item.label}
							component={Link}
							to={item.to}
							active={navItemOpened(item.to)}
						>
							{item.children && (
								<Container mt="sm" ml="md">
									{item.children.map((child) => {
										return (
											<NavLink
												key={child.to}
												label={child.label}
												component={Link}
												to={child.to}
												active={navItemActive(child.to)}
											/>
										);
									})}
								</Container>
							)}
						</NavLink>
					);
				})}
			</AppShell.Navbar>

			<AppShell.Main>
				<Outlet />
			</AppShell.Main>
		</AppShell>
	);
}
