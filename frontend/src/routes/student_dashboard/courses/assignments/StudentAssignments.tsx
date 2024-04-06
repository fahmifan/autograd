import {
	Anchor,
	Box,
	Button,
	Group,
	Pagination,
	Table,
	Title,
	VisuallyHidden,
} from "@mantine/core";
import { notifications } from "@mantine/notifications";
import {
	MDXEditor,
	headingsPlugin,
	listsPlugin,
	markdownShortcutPlugin,
	quotePlugin,
	thematicBreakPlugin,
} from "@mdxeditor/editor";
import { Editor } from "@monaco-editor/react";
import { IconExternalLink, IconFileCheck } from "@tabler/icons-react";
import * as dayjs from "dayjs";
import { useState } from "react";
import {
	Form,
	Link,
	type LoaderFunctionArgs,
	redirect,
	useLoaderData,
	useNavigate,
	useSearchParams,
	useSubmit,
} from "react-router-dom";
import { Breadcrumbs } from "../../../../components/Breadcrumbs";
import type {
	FindAllStudentAssignmentsResponse,
	StudentAssignment,
} from "../../../../pb/autograd/v1/autograd_pb";
import { AutogradServiceClient } from "../../../../service";
import { parseIntWithDefault } from "../../../../types/parser";
import { useCourseDetail } from "../hooks";

export function ListStudentAssignments() {
	const res = useLoaderData() as FindAllStudentAssignmentsResponse;
	const navigate = useNavigate();

	const [searchParams] = useSearchParams();
	const courseID = searchParams.get("courseID") ?? "";

	const hookCourse = useCourseDetail({ courseID })

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

	const items = [
		{ title: "Courses", to: "/student-dashboard/courses" },
		{
			title: hookCourse?.res?.course?.name ?? "",
			to: `/student-dashboard/courses/detail?courseID=${courseID}`,
		},
		{
			title: "Assignments",
			to: `/student-dashboard/courses/assignments?courseID=${courseID}`,
		},
	];

	return (
		<>
			<Breadcrumbs items={items} />

			<Title order={2} mb="lg" mt="lg">
				Assignments
			</Title>
			<Table striped highlightOnHover mb="lg" maw={700}>
				<Table.Thead>
					<Table.Tr>
						<Table.Th>Name</Table.Th>
						<Table.Th>Deadline</Table.Th>
						<Table.Th>Last Update</Table.Th>
						<Table.Th>Submitted</Table.Th>
						<Table.Th>Detail</Table.Th>
					</Table.Tr>
				</Table.Thead>

				<Table.Tbody>
					{res.assignments.map((assg) => {
						return (
							<Table.Tr key={assg.id}>
								<Table.Td>{assg.name}</Table.Td>
								<Table.Td>{humanizeDate(assg.deadlineAt)}</Table.Td>
								<Table.Td>{humanizeDate(assg.updatedAt)}</Table.Td>
								<Table.Td align="center">
									{assg.hasSubmission && <IconFileCheck color="green" />}
								</Table.Td>
								<Table.Td align="center">
									<Anchor
										component={Link}
										to={`/student-dashboard/courses/assignments/detail?courseID=${courseID}&id=${assg.id}`}
									>
										<IconExternalLink color="#339AF0" />
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
					navigate(
						`/student-dashboard/courses/assignments?courseID=${courseID}&page=${page}&limit=${res.paginationMetadata?.limit}`,
					);
				}}
				siblings={1}
				boundaries={2}
			/>
		</>
	);
}

export function DetailStudentAssignment() {
	const res = useLoaderData() as StudentAssignment;
	const submit = useSubmit();
	const [submisisonCode, setSubmissionCode] = useState<string>("");

	const [searchParams] = useSearchParams();
	const courseID = searchParams.get("courseID") ?? "";
	const hookCourse = useCourseDetail({ courseID })

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
		{ title: "Courses", to: "/student-dashboard/courses" },
		{
			title: hookCourse?.res?.course?.name ?? "",
			to: `/student-dashboard/courses/detail?courseID=${courseID}`,
		},
		{ title: "Assignments", to: `/student-dashboard/courses/assignments?courseID=${courseID}` },
		{
			title: res.name,
			to: `/student-dashboard/courses/assignments/detail?courseID=${courseID}&id=${res.id}`,
		},
	];

	return (
		<>
			<Breadcrumbs items={items} />

			<Title order={2} mb="lg" mt="lg">
				Assignments
			</Title>

			<Title order={3} mb="sm" mt="lg">
				{res.name}
			</Title>

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
						<Table.Tr>
							<Table.Th>Submitted At</Table.Th>
							<Table.Td>
								{res.hasSubmission &&
									humanizeDate(res.submission?.updatedAt ?? "-")}
							</Table.Td>
						</Table.Tr>
						<Table.Tr>
							<Table.Th>Grade</Table.Th>
							<Table.Td>
								{res.hasSubmission && res.submission?.isGraded
									? res.submission?.grade
									: "-"}
							</Table.Td>
						</Table.Tr>
					</Table.Tbody>
				</Table>
			</Group>

			<Box mt="md">
				<Title order={4} my="lg">
					Description
				</Title>
				<Box
					p="lg"
					style={{
						border: "1px solid #e0e0e0",
						borderRadius: "8px",
					}}
				>
					<MDXEditor
						markdown={res.description}
						readOnly
						plugins={[
							headingsPlugin(),
							listsPlugin(),
							quotePlugin(),
							thematicBreakPlugin(),
							markdownShortcutPlugin(),
						]}
					/>
				</Box>
			</Box>

			<Box>
				<Group>
					<Title order={4} my="lg">
						Submission
					</Title>
					{res.hasSubmission ? (
						<Form method="post" id="update-student-submission">
							<VisuallyHidden>
								<input type="hidden" name="assignment_id" value={res.id} />
								<input
									type="hidden"
									name="submission_id"
									value={res.submission?.id ?? ""}
								/>
								<input
									type="hidden"
									name="submission_code"
									value={submisisonCode || res.submission?.submissionCode}
								/>
							</VisuallyHidden>
							<Button
								variant="light"
								color="green.9"
								type="submit"
								name="intent"
								value="resubmit_submission"
								onClick={(event) => {
									event.preventDefault();
									const ok = confirm("Are you sure to resubmit?");
									if (!ok) {
										return;
									}
									submit(event.currentTarget);
								}}
							>
								Resubmit
							</Button>
						</Form>
					) : (
						<Form method="post" id="create-student-submission">
							<VisuallyHidden>
								<input type="hidden" name="assignment_id" value={res.id} />
								<input
									type="hidden"
									name="submission_code"
									value={submisisonCode}
								/>
							</VisuallyHidden>

							<Button
								type="submit"
								variant="light"
								size="sm"
								name="intent"
								value="submit_submission"
								onClick={(event) => {
									event.preventDefault();
									const ok = confirm("Are you sure to submit?");
									if (!ok) {
										return;
									}
									submit(event.currentTarget);
								}}
							>
								Submit
							</Button>
						</Form>
					)}
				</Group>

				<Box
					py="lg"
					style={{
						border: "1px solid #e0e0e0",
						borderRadius: "8px",
					}}
				>
					<Editor
						onChange={(value) => {
							setSubmissionCode(value as string);
						}}
						height="300px"
						defaultLanguage="cpp"
						language="cpp"
						defaultValue={
							res?.hasSubmission
								? res.submission?.submissionCode
								: res.codeTemplate
						}
					/>
				</Box>
			</Box>
		</>
	);
}

function humanizeDate(date: string): string {
	if (date === "") {
		return "";
	}
	return dayjs(date).format("HH:MM - ddd DD MMM YYYY");
}

export async function loaderDetailStudentAssignment({
	request,
}: LoaderFunctionArgs): Promise<StudentAssignment> {
	const url = new URL(request.url);
	const id = url.searchParams.get("id");

	const res = await AutogradServiceClient.findStudentAssignment({
		id: id as string,
	});

	return res;
}

export async function loaderListStudentAssignments({
	request,
}: LoaderFunctionArgs): Promise<FindAllStudentAssignmentsResponse> {
	const url = new URL(request.url);
	const page = parseIntWithDefault(url.searchParams.get("page"), 1);
	const limit = parseIntWithDefault(url.searchParams.get("limit"), 10);

	const res = await AutogradServiceClient.findAllStudentAssignments({
		paginationRequest: {
			limit,
			page,
		},
	});

	return res;
}

export async function actionDetailAssignment({
	request,
}: LoaderFunctionArgs): Promise<Response | null> {
	const form = await request.formData();
	const intent = form.get("intent");

	switch (intent) {
		case "submit_submission": {
			const assignmentId = form.get("assignment_id") as string;
			const submissionCode = form.get("submission_code") as string;

			await AutogradServiceClient.submitStudentSubmission({
				assignmentId,
				submissionCode,
			});

			notifications.show({
				title: "Submit",
				message: "Success submit your submission",
				withCloseButton: true,
				color: "green",
			});

			return redirect(
				`/student-dashboard/courses/assignments/detail?id=${assignmentId}`,
			);
		}
		case "resubmit_submission": {
			const assignmentId = form.get("assignment_id") as string;
			const submissionId = form.get("submission_id") as string;
			const submissionCode = form.get("submission_code") as string;

			await AutogradServiceClient.resubmitStudentSubmission({
				submissionId,
				submissionCode,
			});

			notifications.show({
				title: "Resubmit",
				message: "Success resubmit your submission",
				withCloseButton: true,
				color: "green",
			});
			return redirect(
				`/student-dashboard/courses/assignments/detail?id=${assignmentId}`,
			);
		}
		default:
			throw new Error("Invalid intent");
	}
}
