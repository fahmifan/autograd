import {
	ActionIcon,
	Button,
	Card,
	FileInput,
	Flex,
	Input,
	Stack,
	Text,
	TextInput,
	Title,
	Tooltip,
	VisuallyHidden,
	rem,
} from "@mantine/core";
import { DateTimePicker } from "@mantine/dates";
import type { MDXEditorMethods } from "@mdxeditor/editor";
import { Editor } from "@monaco-editor/react";
import { IconTrash, IconUpload } from "@tabler/icons-react";
import { useRef, useState } from "react";
import { useMutation } from "react-query";
import {
	type ActionFunctionArgs,
	Form,
	type LoaderFunctionArgs,
	redirect,
	useLoaderData,
	useSearchParams,
	useSubmit,
} from "react-router-dom";
import { Breadcrumbs } from "../../../../components/Breadcrumbs";
import type { Assignment } from "../../../../pb/autograd/v1/autograd_pb";
import { AutogradRPCClient, AutogradServiceClient } from "../../../../service";
import { MarkdownEditor } from "./shared";

export function DetailAssignment() {
	const res = useLoaderData() as Assignment;

	const [searchParams] = useSearchParams();
	const courseID = searchParams.get("courseID") ?? "";

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
		{ title: "Courses", to: "/backoffice/courses" },
		{
			title: res.course?.name ?? "",
			to: `/backoffice/courses/detail?courseID=${courseID}`,
		},
		{
			title: "Assignments",
			to: `/backoffice/courses/assignments?courseID=${courseID}`,
		},
		{
			title: res.name,
			to: `/backoffice/courses/assignments/detail?courseID=${courseID}&id=${res.id}`,
		},
	];

	return (
		<>
			<Breadcrumbs items={items} />

			<Flex mt="lg" direction="row" justify="space-between">
				<Title order={3} mb="lg">
					{res.name}
				</Title>
				<Form
					method="POST"
					id="delete-assignment"
					onSubmit={(e) => {
						e.preventDefault();

						const ok = confirm(
							`Are you sure you want to delete assignment "${res.name}"?`,
						);
						if (!ok) {
							return;
						}

						submit(e.currentTarget);
					}}
				>
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
							name="courseID"
							id="courseID"
							value={courseID}
						/>

						<Input type="hidden" name="id" id="id" value={res.id} />

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
						placeholder={
							res?.caseInputFile ? res?.caseInputFile?.id : "Select file"
						}
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
						placeholder={
							res?.caseOutputFile ? res?.caseOutputFile?.id : "Select file"
						}
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

export async function actionDetailAssignment(
	arg: ActionFunctionArgs,
): Promise<Response | null> {
	const form = await arg.request.formData();
	const intent = form.get("intent");
	switch (intent) {
		case "delete-assignment":
			return await doDeleteAssignment(form);
		case "update-assignment":
			return await doUpdateAssignment(form);
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

async function doUpdateAssignment(form: FormData): Promise<Response | null> {
	const id = form.get("id") as string;
	const courseID = form.get("courseID") as string;
	const name = form.get("name") as string;
	const description = form.get("description") as string;
	const caseInputFileId = form.get("case_input_file_id") as string;
	const caseOutputFileId = form.get("case_output_file_id") as string;
	const deadlineAt = form.get("deadline_at") as string;
	const template = form.get("template") as string;

	const res = await AutogradServiceClient.updateAssignment({
		id,
		name,
		description,
		caseInputFileId,
		caseOutputFileId,
		deadlineAt,
		template,
	});

	if (res) {
		return redirect(
			`/backoffice/assignments/detail?courseID=${courseID}&id=${id}`,
		);
	}

	return null;
}
async function doDeleteAssignment(form: FormData): Promise<Response | null> {
	const id = form.get("id") as string;
	await AutogradServiceClient.deleteAssignment({
		id,
	});
	return redirect("/backoffice/assignments");
}
