import {
	Button,
	FileInput,
	Group,
	Input,
	Table,
	TextInput,
	Title,
	rem,
} from "@mantine/core";
import { DateTimePicker } from "@mantine/dates";
import { IconUpload } from "@tabler/icons-react";
import { useState } from "react";
import { useMutation } from "react-query";
import { ActionFunctionArgs, Form, useLoaderData } from "react-router-dom";
import { FindAllAssignmentsResponse } from "../../pb/autograd/v1/autograd_pb";
import { AutogradRPCClient, AutogradServiceClient } from "../../service";

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
			<Title order={3}>Assignments</Title>
			<Table striped highlightOnHover>
				<Table.Thead>
					<Table.Tr>
						<Table.Th>ID</Table.Th>
						<Table.Th>Name</Table.Th>
						<Table.Th>Assigner</Table.Th>
					</Table.Tr>
				</Table.Thead>

				<Table.Tbody>
					{res?.assignments?.map((assignment) => {
						return (
							<Table.Tr key={assignment.id}>
								<Table.Td>{assignment.id}</Table.Td>
								<Table.Td>{assignment.name}</Table.Td>
								<Table.Td>{assignment.assigner?.name ?? ""}</Table.Td>
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
		<div>
			<Title order={3}>Create Assignment</Title>
			<Group>
				<Form method="POST" id="create-assignment">
					<TextInput label="Name" required name="name" title="Name" id="name" />
					<TextInput
						required
						label="Description"
						title="Description"
						name="description"
						id="description"
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

					<Button mt="md" type="submit">
						Create
					</Button>
				</Form>
			</Group>
		</div>
	);
}

export async function loaderListAssignments(): Promise<FindAllAssignmentsResponse> {
	return await AutogradServiceClient.findAllAssignments({
		limit: 10,
		page: 1,
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

	await AutogradServiceClient.createAssignment({
		name,
		description,
		caseInputFileId,
		caseOutputFileId,
		deadlineAt,
	});

	return null;
}
