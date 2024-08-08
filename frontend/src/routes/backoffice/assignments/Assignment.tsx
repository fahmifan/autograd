import {
	ActionIcon,
	Anchor,
	Breadcrumbs,
	Button,
	Card,
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
import { Editor } from "@monaco-editor/react";
import { IconExternalLink, IconNote, IconTrash, IconUpload } from "@tabler/icons-react";
import { forwardRef, useRef, useState } from "react";
import { useMutation, useQuery } from "react-query";
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
	PaginationRequest,
} from "../../../pb/autograd/v1/autograd_pb";
import { AutogradRPCCmdClient, AutogradCmdClient, AutogradQueryClient } from "../../../service";

export function ListAssignments() {
	const { data: res } = useQuery({
		queryKey: ["list_assignments"],
		queryFn: () => {
			return AutogradQueryClient.findAllAssignments({
				paginationRequest: new PaginationRequest({
					limit: 10,
					page: 1,
				})
			});
		}
	});
	const submit = useSubmit();

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
						<Table.Th className="text-center">Action</Table.Th>
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
									<Flex direction="row">
										<Anchor
												component={Link}
												to={`/backoffice/assignments/detail?id=${assignment.id}`}
												size="sm"
												mr="sm"
											>
												<Tooltip label={`Detail Assignment for ${assignment.name}`}>
													<IconExternalLink color="#339AF0" />
												</Tooltip>
											</Anchor>
										<Anchor
											component={Link}
											to={`/backoffice/assignments/submissions?assignmentID=${assignment.id}`}
											size="sm"
											mr="sm"
										>
											<Tooltip label={`Submission for ${assignment.name}`}>
												<IconNote color="#339AF0" />
											</Tooltip>
										</Anchor>
										<Form method="POST" id="delete-assignment" onSubmit={e => {
											e.preventDefault();
											const ok = confirm(`Are you sure you want to delete assignment "${assignment.name}"?`);
											if (!ok) {
												return;
											}
											submit(e.currentTarget)
										}}>
											<VisuallyHidden>
												<input name="id" value={assignment.id} />
											</VisuallyHidden>
											<Tooltip label={`Delete assignment ${assignment.name}`}>
												<ActionIcon type="submit" name="intent" value="delete-assignment" variant="outline" aria-label="Delete assignment" color="red.5" size="sm">
													<IconTrash />
												</ActionIcon>
											</Tooltip>
										</Form>
									</Flex>
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
	const [template, setTemplate] = useState("");
	const submit = useSubmit();

	const mutateUploadStdin = useMutation({
		mutationKey: "uploadStdin",
		mutationFn: async (file: File) => {
			const res = await AutogradRPCCmdClient.saveMedia({
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
			const res = await AutogradRPCCmdClient.saveMedia({
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

						<Input 
							type="hidden"
							name="template"
							id="template"
							value={template}
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

				<Text mt="sm">Template</Text>
				<Card shadow="sm" padding="lg" radius="md" withBorder maw={800}>
					<Editor
						onChange={(value) => {
							setTemplate(value as string);
						}}
						height="300px"
						defaultLanguage="cpp"
						language="cpp"
						value={template}
						defaultValue={""}
					/>
				</Card>

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
	const [template, setTemplate] = useState(res.template);
	const markdownRef = useRef<MDXEditorMethods>(null);
	const submit = useSubmit();

	const mutateUploadStdin = useMutation({
		mutationKey: "uploadStdin",
		mutationFn: async (file: File) => {
			const res = await AutogradRPCCmdClient.saveMedia({
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
			const res = await AutogradRPCCmdClient.saveMedia({
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

			<Flex direction="row" justify="space-between">
				<Title order={3} mb="lg">
					{res.name}
				</Title>
				<Form method="POST" id="delete-assignment"
					onSubmit={e => {
						e.preventDefault();

						const ok = confirm(`Are you sure you want to delete assignment "${res.name}"?`);
						if (!ok) {
							return;
						}

						submit(e.currentTarget);
					}}>
					<VisuallyHidden>
						<input name="id" value={res.id} />
					</VisuallyHidden>
					<Tooltip label="Delete Assignment">
						<ActionIcon type="submit" variant="outline" color="red.5" size="md">
							<IconTrash aria-label="Delete assignment" />
						</ActionIcon>
					</Tooltip>
				</Form>
			</Flex>
			<Form method="post" id="update-assignment">
				<Stack maw={400}>
					<TextInput
						label="Name"
						required
						name="name"
						title="Name"
						id="name"
						defaultValue={res.name}
					/>

					<VisuallyHidden>
						<Input
							type="hidden"
							name="id"
							id="id"
							value={res.id}
						/>

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

						<Input 
							type="hidden"
							name="template"
							value={template}
							defaultValue={res.template}
						/>
					</VisuallyHidden>

					<FileInput
						required
						label="Case Input/Stdin"
						title="Case Input/Stdin"
						placeholder={res?.caseInputFile? res?.caseInputFile?.id : "Select file"}
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
						placeholder={res?.caseOutputFile? res?.caseOutputFile?.id : "Select file"}
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
						defaultValue={new Date(res.deadlineAt)}
					/>
				</Stack>

				<Text mt="sm">Template</Text>
				<Card shadow="sm" padding="lg" radius="md" withBorder maw={800}>
					<Editor
						onChange={(value) => {
							setTemplate(value as string);
						}}
						height="300px"
						defaultLanguage="cpp"
						language="cpp"
						value={template}
						defaultValue={res.template ? res.template : `// ${res.name}`}
					/>
				</Card>

				<Text py="lg">Description</Text>

				<MarkdownEditor ref={markdownRef} defaultValue={res.description} />

				<Button
					mt="md"
					type="submit"
					name="intent" 
					value="update-assignment"
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
	return await AutogradCmdClient.findAllAssignments({
		paginationRequest: {
			limit: 10,
			page: 1,
		},
	});
}

export async function actionListAssignments(arg: ActionFunctionArgs,): Promise<Response | null> {
	const formData = await arg.request.formData();
	const id = formData.get("id") as string;

	const res = await AutogradCmdClient.deleteAssignment({
		id,
	})

	if (res) {
		return redirect("/backoffice/assignments");
	}

	return null;
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
	const template = formData.get("template") as string;

	const res = await AutogradCmdClient.createAssignment({
		name,
		description,
		caseInputFileId,
		caseOutputFileId,
		deadlineAt,
		template,
	});

	if (res) {
		return redirect("/backoffice/assignments");
	}

	return null;
}

export async function actionDetailAssignment(arg: ActionFunctionArgs): Promise<Response | null> {
	const form = await arg.request.formData();
	const intent = form.get("intent");
	switch (intent) {
		case "delete-assignment": return await doDeleteAssignment(form);
		case "update-assignment": return await doUpdateAssignment(form);
	}

	return null;
}

async function doDeleteAssignment(form: FormData): Promise<Response | null> {
	const id = form.get("id") as string;
	await AutogradCmdClient.deleteAssignment({
		id,
	});
	return redirect("/backoffice/assignments");
}

async function doUpdateAssignment(form: FormData): Promise<Response | null> {
	const id = form.get("id") as string;
	const name = form.get("name") as string;
	const description = form.get("description") as string;
	const caseInputFileId = form.get("case_input_file_id") as string;
	const caseOutputFileId = form.get("case_output_file_id") as string;
	const deadlineAt = form.get("deadline_at") as string;
	const template = form.get("template") as string;

	const res = await AutogradCmdClient.updateAssignment({
		id,
		name,
		description,
		caseInputFileId,
		caseOutputFileId,
		deadlineAt,
		template
	});

	if (res) {
		return redirect(`/backoffice/assignments/detail?id=${id}`);
	}

	return null;

}

export async function loadEditAssignment({
	request,
}: LoaderFunctionArgs): Promise<Assignment> {
	const url = new URL(request.url);
	const id = url.searchParams.get("id") as string;

	const res = await AutogradCmdClient.findAssignment({
		id,
	});

	return res;
}
