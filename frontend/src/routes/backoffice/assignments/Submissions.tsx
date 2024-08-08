import { Anchor, Box, Breadcrumbs, Table, Text, Title } from "@mantine/core";
import { Editor } from "@monaco-editor/react";
import { IconExternalLink } from "@tabler/icons-react";
import { Link, LoaderFunctionArgs, useLoaderData } from "react-router-dom";
import { FindAllSubmissionsForAssignmentResponse, Submission } from "../../../pb/autograd/v1/autograd_pb";
import { AutogradCmdClient } from "../../../service";

export function ListSubmissions() {	
	const res = useLoaderData() as FindAllSubmissionsForAssignmentResponse;

	const items = [
		{ title: "Assignments", to: "/backoffice/assignments" },
		{ title: "Submission", to: `/backoffice/assignments/submissions?assignmentID=${res.assignmentId}`},
	].map((item) => {
		return (
			<Anchor key={item.to} component={Link} to={item.to}>
				{item.title}
			</Anchor>
		);
	});

	if (!res || !res.submissions || res.submissions.length === 0) {
		return (
			<>
				<Breadcrumbs mb="lg">{items}</Breadcrumbs>
				<p>
					<i>No Submissions</i>
				</p>
			</>
		);
	}


	return (
		<div>
			<Breadcrumbs mb="lg">{items}</Breadcrumbs>

			<Title order={3} mb="lg">
				Submissions
			</Title>
			<Table striped highlightOnHover maw={700} mb="lg">
				<Table.Thead>
					<Table.Tr>
						<Table.Th>ID</Table.Th>
						<Table.Th>Submitter</Table.Th>
						<Table.Th>Detail</Table.Th>
					</Table.Tr>
				</Table.Thead>

				<Table.Tbody>
					{res?.submissions?.map((subm) => {
						return (
							<Table.Tr key={subm.id}>
								<Table.Td>{subm.id}</Table.Td>
								<Table.Td>{subm.submitterName}</Table.Td>
								<Table.Td>
									<Anchor
										component={Link}
										to={`/backoffice/assignments/submissions/detail?submissionID=${subm.id}`}
									>
										<IconExternalLink color="#339AF0" />
									</Anchor>
								</Table.Td>
							</Table.Tr>
						);
					})}
				</Table.Tbody>
			</Table>
		</div>
	);
}

export function SubmissionDetail() {
	const res = useLoaderData() as Submission;

	const items = [
		{ title: "Assignments", to: "/backoffice/assignments" },
		{ title: "Submission", to: `/backoffice/assignments/submissions?assignmentID=${res.assignment?.id}`},
		{ title: res.submitter?.name, to: `/backoffice/assignments/submissions/detail?submissionID=${res.id}`},
	].map((item) => {
		return (
			<Anchor key={item.to} component={Link} to={item.to}>
				{item.title}
			</Anchor>
		);
	});

	return (
		<div>
			<Breadcrumbs mb="lg">{items}</Breadcrumbs>

			<Title order={3} mb="lg">
				Submission
			</Title>

			<Text mb="sm">
				Student: {res.submitter?.name}
			</Text>

			<Box
				py="lg"
				style={{
					border: "1px solid #e0e0e0",
					borderRadius: "8px",
				}}
			>
				<Editor
					height="300px"
					defaultLanguage="cpp"
					language="cpp"
					defaultValue={res?.submissionCode ?? "// some comment"}
				/>
			</Box>
		</div>
	);
}

export async function loaderListSubmissions({
	request,
}: LoaderFunctionArgs): Promise<FindAllSubmissionsForAssignmentResponse> {
	const url = new URL(request.url);
	const assignmentID = url.searchParams.get("assignmentID") as string;

	return await AutogradCmdClient.findAllSubmissionForAssignment({
		assignmentId: assignmentID
	});
}

export async function loaderSubmissionDetail({
	request,
}: LoaderFunctionArgs): Promise<Submission> {
	const url = new URL(request.url);
	const id = url.searchParams.get("submissionID");

	return await AutogradCmdClient.findSubmission({
		id: id as string,
	});
}