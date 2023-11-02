import { AppShell, Container, NavLink, Text } from "@mantine/core";
import { RouteObject } from "react-router-dom";
import { Link, Outlet, useLocation } from "react-router-dom";
import {
	DetailStudentAssignment,
	ListStudentAssignments,
	actionDetailAssignment,
	loaderDetailStudentAssignment,
	loaderListStudentAssignments,
} from "./StudentAssignments";

export const router: RouteObject[] = [
	{
		path: "/student-dashboard",
		element: <DashboardLayout />,
		children: [
			{
				path: "/student-dashboard/assignments",
				element: <ListStudentAssignments />,
				loader: loaderListStudentAssignments,
			},
			{
				path: "/student-dashboard/assignments/detail",
				element: <DetailStudentAssignment />,
				loader: loaderDetailStudentAssignment,
				action: actionDetailAssignment,
			},
		],
	},
];

type NavItem = {
	label: string;
	to: string;
	children: NavItem[];
};

export default function DashboardLayout() {
	const location = useLocation();

	const navitems: NavItem[] = [
		{
			label: "Assignments",
			to: "/student-dashboard/assignments",
			children: [],
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
							{item.children && item.children.length > 0 && (
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
