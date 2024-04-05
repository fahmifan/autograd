import { AppShell, Button, Container, Flex, NavLink, Text } from "@mantine/core";
import { RouteObject } from "react-router-dom";
import { Link, Outlet, useLocation } from "react-router-dom";
import { PrivateRoute } from "../private_route";
import { PageCourseDetail, actionDeleteAssignment } from "./courses/PageCourseDetail";
import { PageCourses } from "./courses/PageCourses";
import { NewAssignment, actionCreateAssignemnt } from "./courses/assignments/CreateAssignment";
import { DetailAssignment, actionDetailAssignment, loadEditAssignment } from "./courses/assignments/DetailAssignment";
import { PageAssignments } from "./courses/assignments/PageAssignments";
import { ListSubmissions, SubmissionDetail, loaderListSubmissions, loaderSubmissionDetail } from "./courses/assignments/Submissions";
import { PageStudents } from "./courses/students/PageStudents";
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
			// User Management
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
			// Courses
			{
				path: "/backoffice/courses",
				element: <PrivateRoute element={<PageCourses />} />,
			},
			{
				path: "/backoffice/courses/detail",
				element: <PrivateRoute element={<PageCourseDetail />} />,
				action: actionDeleteAssignment,
			},
			// Courses Assignments
			{
				path: "/backoffice/courses/assignments",
				element: <PrivateRoute element={<PageAssignments />} />,
				action: actionDeleteAssignment,
			},
			{
				path: "/backoffice/courses/assignments/new",
				element: <PrivateRoute element={<NewAssignment />} />,
				action: actionCreateAssignemnt,
			},
			{
				path: "/backoffice/courses/assignments/detail",
				element: <PrivateRoute element={<DetailAssignment />} />,
				loader: loadEditAssignment,
				action: actionDetailAssignment,
			},
			// Courses Assignments Submission
			{
				path: "/backoffice/courses/assignments/submissions",
				element: <PrivateRoute element={<ListSubmissions />} />,
				loader: loaderListSubmissions,
			},
			{
				path: "/backoffice/courses/assignments/submissions/detail",
				element: <PrivateRoute element={<SubmissionDetail />} />,
				loader: loaderSubmissionDetail,
			},
			// Courses Students
			{
				path: "/backoffice/courses/students",
				element: <PrivateRoute element={<PageStudents />} />
			}
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
			label: "Courses",
			to: "/backoffice/courses"
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
