import { AppShell, Box, Button, Container, Flex, NavLink, Text } from "@mantine/core";
import { RouteObject } from "react-router-dom";
import { Link, Outlet, useLocation } from "react-router-dom";
import { PrivateRoute } from "../private_route";
import {
	CreateAssignment,
	DetailAssignment,
	ListAssignments,
	actionCreateAssignemnt,
	actionDetailAssignment,
	actionListAssignments,
	loadEditAssignment,
	loaderListAssignments,
} from "./assignments/Assignment";
import { ListSubmissions, SubmissionDetail, loaderListSubmissions, loaderSubmissionDetail } from "./assignments/Submissions";
import {
	CreateManagedUser,
	ListManagedUsers,
	actionCreateManagedUser,
	loaderUserManagement,
} from "./user_management/UserManagement";

export const router: RouteObject[] = [
	{
		path: "/backoffice",
		element: <PrivateRoute element={<DashboardLayout />} />,
		children: [
			{
				path: "/backoffice/user-management",
				index: true,
				element: <PrivateRoute element={<ListManagedUsers />} />,
				loader: loaderUserManagement,
			},
			{
				path: "/backoffice/user-management/create",
				element: <PrivateRoute element={<CreateManagedUser />} />,
				action: actionCreateManagedUser,
			},
			{
				path: "/backoffice/assignments",
				element: <PrivateRoute element={<ListAssignments />} />,
				loader: loaderListAssignments,
				action: actionListAssignments,
			},
			{
				path: "/backoffice/assignments/create",
				element: <PrivateRoute element={<CreateAssignment />} />,
				action: actionCreateAssignemnt,
			},
			{
				path: "/backoffice/assignments/detail",
				element: <PrivateRoute element={<DetailAssignment />} />,
				loader: loadEditAssignment,
				action: actionDetailAssignment,
			},
			{
				path: "/backoffice/assignments/submissions",
				element: <PrivateRoute element={<ListSubmissions />} />,
				loader: loaderListSubmissions,
			},
			{
				path: "/backoffice/assignments/submissions/detail",
				element: <PrivateRoute element={<SubmissionDetail />} />,
				loader: loaderSubmissionDetail,
			},
		]
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
		}
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
				<Flex direction="row" justify="space-between" align="center">
					<Text py="sm" px="sm" size="lg">
						Autograd Dashboard
					</Text>
					<Link to="/logout">
						<Button mr="sm" size="compact-sm" color="gray" variant="outline">Logout</Button>
					</Link>
				</Flex>
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
