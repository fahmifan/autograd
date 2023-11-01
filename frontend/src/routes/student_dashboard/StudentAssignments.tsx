import { Anchor, Box, Breadcrumbs, Container, Group, Pagination, Paper, Table, Text, Title } from "@mantine/core";
import { MDXEditor, headingsPlugin, listsPlugin, markdownShortcutPlugin, quotePlugin, thematicBreakPlugin } from "@mdxeditor/editor";
import * as dayjs from 'dayjs'
import { Link, LoaderFunctionArgs, useLoaderData, useNavigate } from "react-router-dom";
import { FindAllStudentAssignmentsResponse, StudentAssignment } from "../../pb/autograd/v1/autograd_pb";
import { AutogradServiceClient } from "../../service";
import { parseIntWithDefault } from "../../types/parser";

export function ListStudentAssignments() {
    const res = useLoaderData() as FindAllStudentAssignmentsResponse;
    const navigate = useNavigate();

    if (!res || res.assignments.length === 0) {
        return (
            <>
                <h2>Assignments</h2>
                <p>
                    <i>No Assignments</i>
                </p>
            </>
        );
    }

    return (
		<>
			<Title order={2} mb="lg">Assignments</Title>
			<Table striped highlightOnHover mb="lg" maw={600}>
				<Table.Thead>
					<Table.Tr>
						<Table.Th>Name</Table.Th>
						<Table.Th>Deadline</Table.Th>
						<Table.Th>Last Update</Table.Th>
						<Table.Th> </Table.Th>
					</Table.Tr>
				</Table.Thead>

				<Table.Tbody>
					{res.assignments.map((assg) => {
						return (
							<Table.Tr key={assg.id}>
								<Table.Td>{assg.name}</Table.Td>
								<Table.Td>{humanizeDate(assg.deadlineAt)}</Table.Td>
								<Table.Td>{humanizeDate(assg.updatedAt)}</Table.Td>
								<Table.Td>
									<Anchor 
										component={Link}
										to={`/student-dashboard/assignments/detail?id=${assg.id}`} 
										>
											Detail
									</Anchor>
								</Table.Td>
							</Table.Tr>
						);
					})}
				</Table.Tbody>
			</Table>

			<Pagination 
				total={res.paginationMetadata?.totalPage as number}
				onChange={(page) => {
					navigate(`/student-dashboard/assignments?page=${page}&limit=${res.paginationMetadata?.limit}`);
				}}
				siblings={1}
				boundaries={2} />
		</>
	);
}

export function DetailStudentAssignment() {
	const res = useLoaderData() as StudentAssignment;

	if (!res) {
		return (
			<>
				<p>
					<i>Assignment Not Found</i>
				</p>
			</>
		);
	}

	const items = [
		{ title: "Assignments", to: "/student-dashboard/assignments" },
		{ title: res.name, to: `/student-dashboard/assignments/detail?id=${res.id}` },
	].map((item) => {
		return <Anchor key={item.to} component={Link} to={item.to}>
			{item.title}
		</Anchor>
	})

	return <>
		<Title order={2} mb="lg">Assignments</Title>

		<Breadcrumbs mb="lg">
			{items}
		</Breadcrumbs>

		<Title order={3} mb="sm">{res.name}</Title>

		<Group>
			<Table maw={300}>
				<Table.Tbody>
					<Table.Tr>
						<Table.Th>Deadline</Table.Th>
						<Table.Td>{humanizeDate(res.deadlineAt)}</Table.Td>
					</Table.Tr>
					<Table.Tr>
						<Table.Th>Updated At</Table.Th>
						<Table.Td>{humanizeDate(res.updatedAt)}</Table.Td>
					</Table.Tr>
				</Table.Tbody>
			</Table>
		</Group>

		<Box mt="md">
			<Title order={4} my="lg">Description</Title>
			<Paper shadow="xs" p="xl">
				<MDXEditor 
					markdown={res.description}
					readOnly
					plugins={[
						headingsPlugin(), 
						listsPlugin(), 
						quotePlugin(), 
						thematicBreakPlugin(),
						markdownShortcutPlugin(),
					]} />
			</Paper>
		</Box>

	</>
}

function humanizeDate(date: string): string {
	return dayjs(date).format("HH:MM - ddd DD MMM YYYY")
}

export async function loaderDetailStudentAssignment({ request }: LoaderFunctionArgs): Promise<StudentAssignment> {
	const url = new URL(request.url)
	const id = url.searchParams.get("id")

	const res = await AutogradServiceClient.findStudentAssignment({
		id: id as string,
	});

	return res;
}

export async function loaderListStudentAssignments({ request }: LoaderFunctionArgs): Promise<FindAllStudentAssignmentsResponse> {
    const url = new URL(request.url)
	const page = parseIntWithDefault(url.searchParams.get("page"), 1)
	const limit = parseIntWithDefault(url.searchParams.get("limit"), 10)

	const res = await AutogradServiceClient.findAllStudentAssignments({
		paginationRequest: {
			limit,
			page,
		}
	});

	return res;
}