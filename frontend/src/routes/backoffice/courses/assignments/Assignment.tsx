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
	useSubmit,
} from "react-router-dom";
import {
	FindAllAssignmentsResponse,
} from "../../../../pb/autograd/v1/autograd_pb";
import { AutogradRPCClient, AutogradServiceClient } from "../../../../service";

function ListAssignments() {
	const res = useLoaderData() as FindAllAssignmentsResponse;
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
												to={`/backoffice/courses/assignments/detail?id=${assignment.id}`}
												size="sm"
												mr="sm"
											>
												<Tooltip label={`Detail Assignment for ${assignment.name}`}>
													<IconExternalLink color="#339AF0" />
												</Tooltip>
											</Anchor>
										<Anchor
											component={Link}
											to={`/backoffice/courses/assignments/submissions?assignmentID=${assignment.id}`}
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
