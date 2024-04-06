import { Anchor, Box, Table, Text, Title } from "@mantine/core";
import { Editor } from "@monaco-editor/react";
import { IconExternalLink } from "@tabler/icons-react";
import {
	Link,
	type LoaderFunctionArgs,
	useLoaderData,
	useSearchParams,
} from "react-router-dom";
import { Breadcrumbs } from "../../../../components/Breadcrumbs";
import type {
	FindAllSubmissionsForAssignmentResponse,
	Submission,
} from "../../../../pb/autograd/v1/autograd_pb";
import { AutogradServiceClient } from "../../../../service";

export function ListSubmissions() {
	const res = useLoaderData() as FindAllSubmissionsForAssignmentResponse;
	const [searchParams] = useSearchParams();
	const courseID = searchParams.get("courseID") ?? "";

	const items = [
		{ title: "Courses", to: `/backoffice/courses/detail?courseID=${courseID}` },
		{
			title: res.course?.name ?? "",
			to: `/backoffice/courses/detail?courseID=${courseID}`,
		},
		{
			title: "Submission",
			to: `/backoffice/courses/assignments/submissions?assignmentID=${res.assignmentId}`,
		},
	];

	if (!res || !res.submissions || res.submissions.length === 0) {
		return (
			<>
				<Box mb="lg">
					<Breadcrumbs items={items} />
				</Box>
				<p>
					<i>No Submissions</i>
				</p>
			</>
		);
	}

	return (
		<div>
			<Breadcrumbs items={items} />

			<Title order={3} my="lg">
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
										to={`/backoffice/courses/assignments/submissions/detail?courseID=${courseID}&submissionID=${subm.id}`}
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

	const [searchParams] = useSearchParams();
	const courseID = searchParams.get("courseID") ?? "";

	const items = [
		{ title: "Courses", to: "/backoffice/courses" },
		{
			title: "Submission",
			to: `/backoffice/courses/assignments/submissions?courseID=${courseID}&assignmentID=${res.assignment?.id}`,
		},
		{
			title: res.submitter?.name ?? "",
			to: `/backoffice/courses/assignments/submissions/detail?courseID=${courseID}&submissionID=${res.id}`,
		},
	];

	return (
		<div>
			<Breadcrumbs items={items} />

			<Title order={3} mb="lg">
				Submission
			</Title>

			<Text mb="sm">Student: {res.submitter?.name}</Text>

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

	return await AutogradServiceClient.findAllSubmissionForAssignment({
		assignmentId: assignmentID,
	});
}

export async function loaderSubmissionDetail({
	request,
}: LoaderFunctionArgs): Promise<Submission> {
	const url = new URL(request.url);
	const id = url.searchParams.get("submissionID");

	return await AutogradServiceClient.findSubmission({
		id: id as string,
	});
}
