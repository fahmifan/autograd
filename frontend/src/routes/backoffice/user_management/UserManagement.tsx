import {
	Button,
	Group,
	Input,
	Pagination,
	Select,
	Table,
	TextInput,
	Title,
} from "@mantine/core";
import {
	ActionFunctionArgs,
	Form,
	LoaderFunctionArgs,
	redirect,
	useLoaderData,
	useNavigate,
	useSubmit,
} from "react-router-dom";
import { FindAllManagedUsersResponse } from "../../../pb/autograd/v1/autograd_pb";
import { AutogradServiceClient } from "../../../service";

export function ListManagedUsers() {
	const { managedUsers, paginationMetadata } =
		useLoaderData() as FindAllManagedUsersResponse;
	const navigate = useNavigate();

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
			<Title order={3} mb="lg">
				List Users
			</Title>
			<Table striped highlightOnHover mb="lg" maw={700}>
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

			<Pagination
				mb="lg"
				total={paginationMetadata?.totalPage as number}
				onChange={(page) => {
					navigate(
						`/backoffice/user-management?page=${page}&limit=${paginationMetadata?.limit}`,
					);
				}}
				siblings={1}
				boundaries={2}
			/>
		</>
	);
}

export function CreateManagedUser() {
	const submit = useSubmit();

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
			<Title order={3} mb="lg">Create User</Title>
			<Group>
				<Form method="post" id="create-managed-user">
						<TextInput mb="md" type="text" name="name" id="name" label="Name" />
						<TextInput mb="md" type="email" name="email" id="email" label="Email" />
						<Select
							mb="md"
							label="Role"
							name="role"
							id="role"
							placeholder="Choose a role"
							data={roleSelection}
						/>

					<Button
						type="submit"
						onClick={(event) => {
							event.preventDefault();
							const ok = confirm("Are you sure you want to create this user?");
							if (!ok) {
								return;
							}
							submit(event.currentTarget);
						}}
					>
						Create User
					</Button>
				</Form>
			</Group>
		</>
	);
}

export async function loaderUserManagement({
	request,
}: LoaderFunctionArgs): Promise<FindAllManagedUsersResponse> {
	const url = new URL(request.url);
	const page = parseIntWithDefault(url.searchParams.get("page"), 1);
	const limit = parseIntWithDefault(url.searchParams.get("limit"), 10);

	const res = await AutogradServiceClient.findAllManagedUsers({
		paginationRequest: {
			limit,
			page,
		},
	});

	return res;
}

function parseIntWithDefault(
	value: string | null,
	defaultValue: number,
): number {
	if (value) {
		return parseInt(value);
	}

	return defaultValue;
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
		return redirect("/backoffice/user-management");
	}

	return null;
}
