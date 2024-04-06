import { Loader, Pagination, Table, Text } from "@mantine/core";
import { useState } from "react";
import { useQuery } from "react-query";
import { useSearchParams } from "react-router-dom";
import { Breadcrumbs } from "../../../../components/Breadcrumbs";
import {
	type FindAllCourseStudentsResponse,
	PaginationRequest,
} from "../../../../pb/autograd/v1/autograd_pb";
import { AutogradServiceClient } from "../../../../service";
import { useAdminCourseDetail } from "../hooks";

function useCourseStudents(arg: {
	courseID: string;
	page: number;
	limit: number;
}): {
	error: unknown;
	res?: FindAllCourseStudentsResponse;
	isLoading: boolean;
} {
	const queryKeys = [
		"courses",
		arg.courseID,
		"students",
		"page",
		arg.page,
		"limit",
		arg.limit,
	];

	const { isLoading, data, isError, error } = useQuery({
		queryKey: queryKeys,
		queryFn: async () => {
			return AutogradServiceClient.findAllCourseStudents({
				courseId: arg.courseID,
				paginationRequest: new PaginationRequest({
					page: arg.page,
					limit: arg.limit,
				}),
			});
		},
	});

	return {
		error,
		res: data,
		isLoading,
	};
}

export function PageStudents() {
	const [searchParams, setSearchParams] = useSearchParams();
	const courseID = searchParams.get("courseID") ?? "";
	const [page, setPage] = useState(
		Number.parseInt(searchParams.get("page") || "1"),
	);
	const limit = Number.parseInt(searchParams.get("limit") || "10");

	const hookCourse = useAdminCourseDetail({ courseID });
	const hookCourseStudents = useCourseStudents({
		courseID,
		page,
		limit,
	});

	const items = [
		{ title: "Courses", to: "/backoffice/courses" },
		{
			title: hookCourse.res?.course?.name ?? "",
			to: `/backoffice/courses/detail?courseID=${courseID}`,
		},
		{
			title: "Students",
			to: `/backoffice/courses/students?courseID=${courseID}`,
		},
	];

	return (
		<section>
			<Breadcrumbs items={items} />

			<Table striped highlightOnHover maw={800} mb="lg">
				<Table.Thead>
					<Table.Tr>
						<Table.Th>Name</Table.Th>
					</Table.Tr>
				</Table.Thead>

				<Table.Tbody>
					{hookCourseStudents.isLoading ? (
						<Loader color="blue" m="lg" size="sm" />
					) : (
						hookCourseStudents.res?.students?.map((student) => {
							return (
								<Table.Tr key={student.id}>
									<Table.Td>
										<Text>{student.name}</Text>
									</Table.Td>
								</Table.Tr>
							);
						})
					)}
				</Table.Tbody>
			</Table>

			<Pagination
				mb="lg"
				total={hookCourseStudents?.res?.paginationMetadata?.totalPage as number}
				value={page}
				onChange={(page) => {
					setPage(page);
					searchParams.set("page", page.toString());
					setSearchParams(searchParams);
				}}
				siblings={1}
				boundaries={2}
			/>
		</section>
	);
}
