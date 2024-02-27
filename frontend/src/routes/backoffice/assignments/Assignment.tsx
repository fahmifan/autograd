import {
	Anchor,
	Breadcrumbs,
	Button,
	FileInput,
	Flex,
	Input,
	Paper,
	Stack,
	Table,
	Text,
	TextInput,
	Title,
	Tooltip,
	VisuallyHidden,
	rem,
} from "@mantine/core";
import { DateTimePicker } from "@mantine/dates";
import {
	MDXEditor,
	MDXEditorMethods,
	headingsPlugin,
	listsPlugin,
	markdownShortcutPlugin,
	quotePlugin,
	thematicBreakPlugin,
} from "@mdxeditor/editor";
import "@mdxeditor/editor/style.css";
import { IconExternalLink, IconNote, IconUpload } from "@tabler/icons-react";
import { forwardRef, useRef, useState } from "react";
import { useMutation } from "react-query";
import {
	ActionFunctionArgs,
	Form,
	Link,
	LoaderFunctionArgs,
	redirect,
	useLoaderData,
	useSubmit,
} from "react-router-dom";
import {
	Assignment,
	FindAllAssignmentsResponse,
} from "../../../pb/autograd/v1/autograd_pb";
import { AutogradRPCClient, AutogradServiceClient } from "../../../service";

export function ListAssignments() {
	const res = useLoaderData() as FindAllAssignmentsResponse;

	if (!res || res.assignments.length === 0) {
		return (
			<>
				<p>
					<i>No Assignments</i>
				</p>
			</>
		);
	}

	return (
		<div>
			<Title order={3} mb="lg">
				Assignments
			</Title>
			<Table striped highlightOnHover maw={800} mb="lg">
				<Table.Thead>
					<Table.Tr>
						<Table.Th>ID</Table.Th>
						<Table.Th>Name</Table.Th>
						<Table.Th>Assigner</Table.Th>
						<Table.Th>Detail</Table.Th>
						<Table.Th>Submissions</Table.Th>
					</Table.Tr>
				</Table.Thead>

				<Table.Tbody>
					{res?.assignments?.map((assignment) => {
						return (
							<Table.Tr key={assignment.id}>
								<Table.Td>{assignment.id}</Table.Td>
								<Table.Td>{assignment.name}</Table.Td>
								<Table.Td>{assignment.assigner?.name ?? ""}</Table.Td>
								<Table.Td>
										<Anchor
											component={Link}
											to={`/backoffice/assignments/detail?id=${assignment.id}`}
											size="sm"
										>
											<Tooltip label={`Detail Assignment for ${assignment.name}`}>
												<IconExternalLink color="#339AF0" />
											</Tooltip>
										</Anchor>
								</Table.Td>
								<Table.Td>
									<Anchor
										component={Link}
										to={`/backoffice/assignments/submissions?assignmentID=${assignment.id}`}
										size="sm"
									>
										<Tooltip label={`Submission for ${assignment.name}`}>
											<IconNote color="#339AF0" />
										</Tooltip>
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

export function CreateAssignment() {
	const [stdinFileID, setStdinFileID] = useState("");
	const [stdoutFileID, setStdoutFileID] = useState("");
	const markdownRef = useRef<MDXEditorMethods>(null);
	const submit = useSubmit();

	const mutateUploadStdin = useMutation({
		mutationKey: "uploadStdin",
		mutationFn: async (file: File) => {
			const res = await AutogradRPCClient.saveMedia({
				file,
				mediaType: "assignment_case_input",
			});
			if (!res.ok) {
				throw new Error("Failed to upload file");
			}

			setStdinFileID(res.value.id ?? "");
			return res;
		},
	});

	const mutateUploadStdout = useMutation({
		mutationKey: "uploadStdout",
		mutationFn: async (file: File) => {
			const res = await AutogradRPCClient.saveMedia({
				file,
				mediaType: "assignment_case_output",
			});
			if (!res.ok) {
				throw new Error("Failed to upload file");
			}

			setStdoutFileID(res.value.id ?? "");
			return res;
		},
	});

	return (
		<>
			<Title order={3}>Create Assignment</Title>
			<Form method="POST" id="create-assignment">
				<Stack maw={400}>
					<TextInput label="Name" required name="name" title="Name" id="name" />

					<VisuallyHidden>
						<Input
							type="hidden"
							name="description"
							id="description"
							value={markdownRef.current?.getMarkdown() ?? ""}
						/>

						<Input
							type="hidden"
							name="case_input_file_id"
							id="case_input_file_id"
							value={stdinFileID}
						/>

						<Input
							type="hidden"
							name="case_output_file_id"
							id="case_output_file_id"
							value={stdoutFileID}
						/>
					</VisuallyHidden>

					<FileInput
						required
						label="Case Input/Stdin"
						title="Case Input/Stdin"
						placeholder="Select file"
						rightSection={
							<IconUpload
								style={{ width: rem(18), height: rem(18) }}
								stroke={1.5}
							/>
						}
						onChange={(event) => {
							if (!event) {
								return;
							}
							mutateUploadStdin.mutateAsync(event);
						}}
					/>

					<FileInput
						required
						label="Case Output/Stdout"
						title="Case Output/Stdout"
						placeholder="Select file"
						rightSection={
							<IconUpload
								style={{ width: rem(18), height: rem(18) }}
								stroke={1.5}
							/>
						}
						onChange={(event) => {
							if (!event) {
								return;
							}
							mutateUploadStdout.mutateAsync(event);
						}}
					/>

					<DateTimePicker
						label="Deadline"
						placeholder="Pick deadline date & time"
						required
						name="deadline_at"
						id="deadline_at"
					/>
				</Stack>

				<Text py="lg">Description</Text>
				<MarkdownEditor ref={markdownRef} />

				<Button
					mt="md"
					type="submit"
					onClick={(event) => {
						event.preventDefault();
						const el = event.currentTarget.form?.elements.namedItem(
							"description",
						) as Element;
						el.setAttribute("value", markdownRef.current?.getMarkdown() ?? "");
						submit(event.currentTarget);
					}}
				>
					Create
				</Button>
			</Form>
		</>
	);
}

export function DetailAssignment() {
	const res = useLoaderData() as Assignment;

	const [stdinFileID, setStdinFileID] = useState(res.caseInputFile?.id ?? "");
	const [stdoutFileID, setStdoutFileID] = useState(
		res.caseOutputFile?.id ?? "",
	);
	const markdownRef = useRef<MDXEditorMethods>(null);
	const submit = useSubmit();

	const mutateUploadStdin = useMutation({
		mutationKey: "uploadStdin",
		mutationFn: async (file: File) => {
			const res = await AutogradRPCClient.saveMedia({
				file,
				mediaType: "assignment_case_input",
			});
			if (!res.ok) {
				throw new Error("Failed to upload file");
			}

			setStdinFileID(res.value.id ?? "");
			return res;
		},
	});

	const mutateUploadStdout = useMutation({
		mutationKey: "uploadStdout",
		mutationFn: async (file: File) => {
			const res = await AutogradRPCClient.saveMedia({
				file,
				mediaType: "assignment_case_output",
			});
			if (!res.ok) {
				throw new Error("Failed to upload file");
			}

			setStdoutFileID(res.value.id ?? "");
			return res;
		},
	});

	const items = [
		{ title: "Assignments", to: "/backoffice/assignments" },
		{ title: res.name, to: `/backoffice/assignments/detail?id=${res.id}` },
	].map((item) => {
		return (
			<Anchor key={item.to} component={Link} to={item.to}>
				{item.title}
			</Anchor>
		);
	});

	return (
		<>
			<Breadcrumbs mb="lg">{items}</Breadcrumbs>

			<Title order={3} mb="lg">
				{res.name}
			</Title>
			<Form method="post" id="edit-assignment">
				<Stack maw={400}>
					<TextInput
						label="Name"
						required
						name="name"
						title="Name"
						id="name"
						value={res.name}
					/>

					<VisuallyHidden>
						<Input
							type="hidden"
							name="description"
							id="description"
							value={markdownRef.current?.getMarkdown() ?? ""}
						/>

						<Input
							type="hidden"
							name="case_input_file_id"
							id="case_input_file_id"
							value={stdinFileID}
						/>

						<Input
							type="hidden"
							name="case_output_file_id"
							id="case_output_file_id"
							value={stdoutFileID}
						/>
					</VisuallyHidden>

					<FileInput
						required
						label="Case Input/Stdin"
						title="Case Input/Stdin"
						placeholder="Select file"
						rightSection={
							<IconUpload
								style={{ width: rem(18), height: rem(18) }}
								stroke={1.5}
							/>
						}
						onChange={(event) => {
							if (!event) {
								return;
							}
							mutateUploadStdin.mutateAsync(event);
						}}
					/>

					<FileInput
						required
						label="Case Output/Stdout"
						title="Case Output/Stdout"
						placeholder="Select file"
						rightSection={
							<IconUpload
								style={{ width: rem(18), height: rem(18) }}
								stroke={1.5}
							/>
						}
						onChange={(event) => {
							if (!event) {
								return;
							}
							mutateUploadStdout.mutateAsync(event);
						}}
					/>

					<DateTimePicker
						label="Deadline"
						placeholder="Pick deadline date & time"
						required
						name="deadline_at"
						id="deadline_at"
						value={new Date(res.deadlineAt)}
					/>
				</Stack>

				<Text py="lg">Description</Text>

				<MarkdownEditor ref={markdownRef} defaultValue={res.description} />

				<Button
					mt="md"
					type="submit"
					onClick={(event) => {
						event.preventDefault();
						const ok = confirm(
							"Are you sure you want to update this assignment?",
						);
						if (!ok) {
							return;
						}

						const el = event.currentTarget.form?.elements.namedItem(
							"description",
						) as Element;
						el.setAttribute("value", markdownRef.current?.getMarkdown() ?? "");
						submit(event.currentTarget);
					}}
				>
					Update
				</Button>
			</Form>
		</>
	);
}

type MarkdownEditorProps = {
	defaultValue?: string;
	onChange?: (value: string) => void;
	ref: React.RefObject<MDXEditorMethods>;
};

export const MarkdownEditor = forwardRef<MDXEditorMethods, MarkdownEditorProps>(
	(props: MarkdownEditorProps, ref) => {
		return (
			<Paper shadow="xs" p="xl">
				<MDXEditor
					ref={ref}
					onChange={props.onChange}
					markdown={props.defaultValue || "## Description"}
					plugins={[
						headingsPlugin(),
						listsPlugin(),
						quotePlugin(),
						thematicBreakPlugin(),
						markdownShortcutPlugin(),
					]}
				/>
			</Paper>
		);
	},
);

export async function loaderListAssignments(): Promise<FindAllAssignmentsResponse> {
	return await AutogradServiceClient.findAllAssignments({
		paginationRequest: {
			limit: 10,
			page: 1,
		},
	});
}

export async function actionCreateAssignemnt(
	arg: ActionFunctionArgs,
): Promise<Response | null> {
	const formData = await arg.request.formData();
	const name = formData.get("name") as string;
	const description = formData.get("description") as string;
	const caseInputFileId = formData.get("case_input_file_id") as string;
	const caseOutputFileId = formData.get("case_output_file_id") as string;
	const deadlineAt = formData.get("deadline_at") as string;

	const res = await AutogradServiceClient.createAssignment({
		name,
		description,
		caseInputFileId,
		caseOutputFileId,
		deadlineAt,
	});

	if (res) {
		return redirect("/backoffice/assignments");
	}

	return null;
}

export async function loadEditAssignment({
	request,
}: LoaderFunctionArgs): Promise<Assignment> {
	const url = new URL(request.url);
	const id = url.searchParams.get("id") as string;

	const res = await AutogradServiceClient.findAssignment({
		id,
	});

	return res;
}
