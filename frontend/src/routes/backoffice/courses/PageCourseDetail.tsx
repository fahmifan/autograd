import {
	ActionIcon,
	Anchor,
	Button,
	Card,
	Flex,
	Grid,
	Pagination,
	Table,
	Text,
	Title,
	Tooltip,
	VisuallyHidden,
} from "@mantine/core";
import "@mdxeditor/editor/style.css";
import { IconExternalLink, IconNote, IconTrash } from "@tabler/icons-react";
import { useState } from "react";
import { useQuery } from "react-query";
import {
	ActionFunctionArgs,
	Form,
	Link,
	redirect,
	useNavigate,
	useSearchParams,
	useSubmit,
} from "react-router-dom";
import { Breadcrumbs } from "../../../components/Breadcrumbs";
import {
	FindAllAssignmentsResponse,
	PaginationMetadata,
} from "../../../pb/autograd/v1/autograd_pb";
import { AutogradServiceClient } from "../../../service";
import { useAdminCourseDetail } from "./hooks";

function useListAssignments(arg: {
	courseId: string;
	page: number;
	limit: number;
}): {
	error: unknown;
	res?: FindAllAssignmentsResponse;
	paginationMetadata?: PaginationMetadata;
} {
	const queryKeys = [
		"courses",
		arg.courseId,
		"assignments",
		arg.page,
		arg.limit,
	];

	const { isLoading, data, isError, error } = useQuery({
		queryKey: queryKeys,
		queryFn: async () => {
			return AutogradServiceClient.findAllAssignments({
				courseId: arg.courseId,
				paginationRequest: {
					page: arg.page,
					limit: arg.limit,
				},
			});
		},
	});

	return {
		error,
		paginationMetadata: data?.paginationMetadata,
		res: data,
	};
}

export function PageCourseDetail() {
	const [searchParams] = useSearchParams();
	const [page] = useState(parseInt(searchParams.get("page") || "1"));
	const limit = parseInt(searchParams.get("limit") || "10");
	const courseID = searchParams.get("courseID") ?? "";

	const hookCourse = useAdminCourseDetail({ courseID });
	const { res } = hookCourse;

	const items = [
		{ title: "Courses", to: "/backoffice/courses" },
		{
			title: res?.course?.name ?? "",
			to: `/backoffice/courses/detail?courseID=${courseID}`,
		},
	];

	return (
		<section>
			<Breadcrumbs items={items} />
			<Title order={3} mt="lg">
				{res?.course?.name}
			</Title>

			<Grid>
				<Grid.Col span={4}>
					<Card
						shadow="sm"
						p="xl"
						component={Link}
						to={`/backoffice/courses/assignments?courseID=${courseID}`}
						m="md"
						style={{
							"&:hover": {
								cursor: "pointer",
							},
						}}
					>
						<Text fw={500} size="xl" mt="md">
							Assignments
						</Text>
					</Card>
				</Grid.Col>

				<Grid.Col span={4}>
					<Card
						shadow="sm"
						p="xl"
						component={Link}
						to={`/backoffice/courses/students?courseID=${courseID}`}
						m="md"
						style={{
							"&:hover": {
								cursor: "pointer",
							},
						}}
					>
						<Text fw={500} size="xl" mt="md">
							Students
						</Text>
					</Card>
				</Grid.Col>
			</Grid>
		</section>
	);
}

export async function actionDeleteAssignment(
	arg: ActionFunctionArgs,
): Promise<Response | null> {
	const formData = await arg.request.formData();
	const id = formData.get("id") as string;

	const res = await AutogradServiceClient.deleteAssignment({
		id,
	});

	if (res) {
		return redirect("/backoffice/assignments");
	}

	return null;
}
