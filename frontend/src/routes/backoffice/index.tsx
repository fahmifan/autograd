import { AppShell, Container, NavLink, Text } from "@mantine/core";
import { RouteObject } from "react-router-dom";
import { Link, Outlet, useLocation } from "react-router-dom";
import {
	CreateAssignment,
	ListAssignments,
	actionCreateAssignemnt,
	loaderListAssignments,
} from "./AssignmentSubmission";
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
				path: "/backoffice/assignment-submission",
				element: <ListAssignments />,
				loader: loaderListAssignments,
			},
			{
				path: "/backoffice/assignment-submission/create",
				element: <CreateAssignment />,
				action: actionCreateAssignemnt,
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
			to: "/backoffice/assignment-submission",
			children: [
				{
					label: "List Assignments",
					to: "/backoffice/assignment-submission",
				},
				{
					label: "Create Assignment",
					to: "/backoffice/assignment-submission/create",
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
