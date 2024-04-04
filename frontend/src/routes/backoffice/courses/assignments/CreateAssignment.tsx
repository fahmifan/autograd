import {
	ActionIcon,
	Anchor,
	Box,
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
import { useMutation } from "react-query";
import {
	ActionFunctionArgs,
	Form,
	Link,
	redirect,
	useLoaderData,
	useSearchParams,
	useSubmit,
} from "react-router-dom";
import {
	FindAllAssignmentsResponse,
} from "../../../../pb/autograd/v1/autograd_pb";
import { AutogradRPCClient, AutogradServiceClient } from "../../../../service";
import { MarkdownEditor } from "./Assignment";


export function NewAssignment() {
	const [stdinFileID, setStdinFileID] = useState("");
	const [stdoutFileID, setStdoutFileID] = useState("");
	const markdownRef = useRef<MDXEditorMethods>(null);
	const [template, setTemplate] = useState("");
	const submit = useSubmit();

	const [searchParams] = useSearchParams()
	const courseID = searchParams.get('courseID') ?? ''

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
							name="courseID"
							id="courseID"
							value={courseID}
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

export async function actionCreateAssignemnt(
	arg: ActionFunctionArgs,
): Promise<Response | null> {
	const formData = await arg.request.formData();
	
	const courseId = formData.get("courseID") as string;
	const name = formData.get("name") as string;
	const description = formData.get("description") as string;
	const caseInputFileId = formData.get("case_input_file_id") as string;
	const caseOutputFileId = formData.get("case_output_file_id") as string;
	const deadlineAt = formData.get("deadline_at") as string;
	const template = formData.get("template") as string;

	const res = await AutogradServiceClient.createAssignment({
		courseId,
		name,
		description,
		caseInputFileId,
		caseOutputFileId,
		deadlineAt,
		template,
	});

	if (res) {
		return redirect(`/backoffice/courses/detail?courseID=${courseId}`);
	}

	return null;
}
