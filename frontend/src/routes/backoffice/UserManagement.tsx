import { Button, Group, Input, Select, Table, Title } from "@mantine/core";
import {
	ActionFunctionArgs,
	Form,
	redirect,
	useLoaderData,
} from "react-router-dom";
import { ManagedUser } from "../../pb/autograd/v1/autograd_pb";
import { AutogradServiceClient } from "../../service";

export function ListManagedUsers() {
	const { managedUsers } = useLoaderData() as LoaderResponse;

	if (!managedUsers || managedUsers.length === 0) {
		return (
			<>
				<h2>List Users</h2>
				<p>
					<i>No Users</i>
				</p>
			</>
		);
	}

	return (
		<>
			<Title order={3}>List Users</Title>
			<Table striped highlightOnHover>
				<Table.Thead>
					<Table.Tr>
						<Table.Th>Name</Table.Th>
						<Table.Th>Email</Table.Th>
						<Table.Th>Role</Table.Th>
					</Table.Tr>
				</Table.Thead>

				<Table.Tbody>
					{managedUsers.map((user) => {
						return (
							<Table.Tr key={user.id}>
								<Table.Td>{user.name}</Table.Td>
								<Table.Td>{user.email}</Table.Td>
								<Table.Td>{user.role}</Table.Td>
							</Table.Tr>
						);
					})}
				</Table.Tbody>
			</Table>
		</>
	);
}

type LoaderResponse = {
	managedUsers: ManagedUser[];
};

export async function loaderUserManagement(): Promise<LoaderResponse> {
	const res = await AutogradServiceClient.findAllManagedUsers({
		limit: 10,
		page: 1,
	});

	return res;
}

export function CreateManagedUser() {
	const roleSelection = [
		{
			value: "admin",
			label: "Admin",
		},
		{
			value: "student",
			label: "Student",
		},
	];

	return (
		<>
			<Title order={3}>Create User</Title>
			<Group>
				<Form method="post" id="create-managed-user">
					<p>
						<label htmlFor="name">Name</label>
						<Input type="text" name="name" id="name" />
					</p>
					<p>
						<label htmlFor="email">Email</label>
						<Input type="email" name="email" id="email" />
					</p>
					<p>
						<label htmlFor="role">Role</label>
						<Select
							name="role"
							id="role"
							placeholder="Choose a role"
							data={roleSelection}
						/>
					</p>

					<Button type="submit">Create User</Button>
				</Form>
			</Group>
		</>
	);
}

export async function actionCreateManagedUser({
	request,
}: ActionFunctionArgs): Promise<Response | null> {
	const formData = await request.formData();

	const name = formData.get("name") as string;
	const email = formData.get("email") as string;
	const role = formData.get("role") as string;

	const res = await AutogradServiceClient.createManagedUser({
		email,
		name,
		role,
	});

	if (res) {
		return redirect("/user-management");
	}

	return null;
}