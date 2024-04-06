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
import { useCourseDetail } from "./hooks";

export function PageCourseDetail() {
    const [searchParams] = useSearchParams()
    const [page] = useState(parseInt(searchParams.get('page') || '1'))
    const limit = parseInt(searchParams.get('limit') || '10')
    const courseID = searchParams.get('courseID') ?? ''

    const hookCourse = useCourseDetail({ courseID })
	const { res } = hookCourse;

	const items = [
		{ title: "Courses", to: "/student-dashboard/courses" },
		{ title: res?.course?.name ?? '', to: `/student-dashboard/courses/detail?courseID=${courseID}` },
	]

	return (
		<section>
			<Breadcrumbs items={items} />
			<Title order={3} mt="lg">{res?.course?.name}</Title>
			
			<Grid>
				<Grid.Col span={4}>
					<Card
						shadow="sm"
						p="xl"
						component={Link}
						to={`/student-dashboard/courses/assignments?courseID=${courseID}`}
						m="md"
						style={{
							'&:hover': {
								cursor: 'pointer'
							}
						}}
					>
						<Text fw={500} size="xl" mt="md">Assignments</Text>
					</Card>
				</Grid.Col>
			</Grid>
		</section>
	);
}
