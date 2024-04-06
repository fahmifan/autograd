import { Card, Grid, Text, Title } from "@mantine/core";
import "@mdxeditor/editor/style.css";
import { Link, useSearchParams } from "react-router-dom";
import { Breadcrumbs } from "../../../components/Breadcrumbs";
import { useCourseDetail } from "./hooks";

export function PageCourseDetail() {
	const [searchParams] = useSearchParams();
	const courseID = searchParams.get("courseID") ?? "";

	const hookCourse = useCourseDetail({ courseID });
	const { res } = hookCourse;

	const items = [
		{ title: "Courses", to: "/student-dashboard/courses" },
		{
			title: res?.course?.name ?? "",
			to: `/student-dashboard/courses/detail?courseID=${courseID}`,
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
						to={`/student-dashboard/courses/assignments?courseID=${courseID}`}
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
			</Grid>
		</section>
	);
}
