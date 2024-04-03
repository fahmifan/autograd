import {
	ActionIcon,
	Anchor,
	Flex,
	Pagination,
	Table,
	Title,
	Tooltip,
	VisuallyHidden,
} from "@mantine/core";
import "@mdxeditor/editor/style.css";
import { IconExternalLink, IconNote, IconTrash } from "@tabler/icons-react";
import { useState } from "react";
import { useQuery } from "react-query";
import {
	Form,
	Link,
	useSearchParams,
	useSubmit,
} from "react-router-dom";
import {
	FindAllAssignmentsResponse,
    PaginationMetadata,
} from "../../../pb/autograd/v1/autograd_pb";
import { AutogradServiceClient } from "../../../service";

function useListAssignments(arg: {
    courseId: string;
    page: number;
    limit: number;
}): {
    error: unknown;
    res?: FindAllAssignmentsResponse;
    paginationMetadata?: PaginationMetadata;
} {
    const queryKeys = ["courses", arg.courseId, "assignments", arg.page, arg.limit]

    const { isLoading, data, isError, error } = useQuery({
        queryKey: queryKeys,
        queryFn: async () => {
            return AutogradServiceClient.findAllAssignments({
                courseId: arg.courseId,
                paginationRequest: {
                    page: arg.page,
                    limit: arg.limit,
                }
            })
        },
    })

    return {
        error,
        paginationMetadata: data?.paginationMetadata,
        res: data
    }
}

export function PageCourseDetail() {
    const [searchParams] = useSearchParams()
    const [page, setPage] = useState(parseInt(searchParams.get('page') || '1'))
    const limit = parseInt(searchParams.get('limit') || '10')
    const courseID = searchParams.get('courseID') ?? ''

    const hookListAssignment = useListAssignments({
        courseId: courseID,
        limit,
        page,
    })

	const { res } = hookListAssignment;
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
		<section>
			<Title order={2} mb="lg">
				Course Assignments
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
												to={`/backoffice/assignments/detail?courseID=${courseID}&id=${assignment.id}`}
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

            <Pagination
                mb="lg"
                total={res.paginationMetadata?.totalPage as number}
                value={page}
                onChange={setPage}
                siblings={1}
                boundaries={2}
            />
		</section>
	);
}
