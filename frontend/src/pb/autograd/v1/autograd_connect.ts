// @generated by protoc-gen-connect-es v0.12.0 with parameter "target=ts"
// @generated from file autograd/v1/autograd.proto (package autograd.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import {
	ActivateManagedUserRequest,
	Assignment,
	CreateAdminCourseRequest,
	CreateAssignmentRequest,
	CreatedResponse,
	CreateManagedUserRequest,
	CreateSubmissionRequest,
	DeleteByIDRequest,
	Empty,
	FindAdminCourseDetailResponse,
	FindAllAdminCoursesResponse,
	FindAllAssignmentsRequest,
	FindAllAssignmentsResponse,
	FindAllCourseStudentsRequest,
	FindAllCourseStudentsResponse,
	FindAllManagedUsersResponse,
	FindAllPaginationRequest,
	FindAllStudentAssignmentsRequest,
	FindAllStudentAssignmentsResponse,
	FindAllStudentEnrolledCoursesRequest,
	FindAllStudentEnrolledCoursesResponse,
	FindAllSubmissionsForAssignmentRequest,
	FindAllSubmissionsForAssignmentResponse,
	FindByIDRequest,
	FindStudentCourseDetailResponse,
	LoginRequest,
	LoginResponse,
	PingResponse,
	ResubmitStudentSubmissionRequest,
	StudentAssignment,
	Submission,
	SubmitStudentSubmissionRequest,
	UpdateAdminCourseRequest,
	UpdateAssignmentRequest,
	UpdateSubmissionRequest,
} from "./autograd_pb.js";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * @generated from service autograd.v1.AutogradService
 */
export const AutogradService = {
	typeName: "autograd.v1.AutogradService",
	methods: {
		/**
		 * @generated from rpc autograd.v1.AutogradService.Ping
		 */
		ping: {
			name: "Ping",
			I: Empty,
			O: PingResponse,
			kind: MethodKind.Unary,
		},
		/**
		 * User Management
		 *
		 * @generated from rpc autograd.v1.AutogradService.CreateManagedUser
		 */
		createManagedUser: {
			name: "CreateManagedUser",
			I: CreateManagedUserRequest,
			O: CreatedResponse,
			kind: MethodKind.Unary,
		},
		/**
		 * @generated from rpc autograd.v1.AutogradService.ActivateManagedUser
		 */
		activateManagedUser: {
			name: "ActivateManagedUser",
			I: ActivateManagedUserRequest,
			O: Empty,
			kind: MethodKind.Unary,
		},
		/**
		 * @generated from rpc autograd.v1.AutogradService.FindAllManagedUsers
		 */
		findAllManagedUsers: {
			name: "FindAllManagedUsers",
			I: FindAllPaginationRequest,
			O: FindAllManagedUsersResponse,
			kind: MethodKind.Unary,
		},
		/**
		 * Assignment Submission
		 * Assignment Queries
		 *
		 * @generated from rpc autograd.v1.AutogradService.FindAssignment
		 */
		findAssignment: {
			name: "FindAssignment",
			I: FindByIDRequest,
			O: Assignment,
			kind: MethodKind.Unary,
		},
		/**
		 * @generated from rpc autograd.v1.AutogradService.FindAllAssignments
		 */
		findAllAssignments: {
			name: "FindAllAssignments",
			I: FindAllAssignmentsRequest,
			O: FindAllAssignmentsResponse,
			kind: MethodKind.Unary,
		},
		/**
		 * @generated from rpc autograd.v1.AutogradService.FindSubmission
		 */
		findSubmission: {
			name: "FindSubmission",
			I: FindByIDRequest,
			O: Submission,
			kind: MethodKind.Unary,
		},
		/**
		 * @generated from rpc autograd.v1.AutogradService.FindAllSubmissionForAssignment
		 */
		findAllSubmissionForAssignment: {
			name: "FindAllSubmissionForAssignment",
			I: FindAllSubmissionsForAssignmentRequest,
			O: FindAllSubmissionsForAssignmentResponse,
			kind: MethodKind.Unary,
		},
		/**
		 * Assignment Command
		 *
		 * @generated from rpc autograd.v1.AutogradService.CreateAssignment
		 */
		createAssignment: {
			name: "CreateAssignment",
			I: CreateAssignmentRequest,
			O: CreatedResponse,
			kind: MethodKind.Unary,
		},
		/**
		 * @generated from rpc autograd.v1.AutogradService.UpdateAssignment
		 */
		updateAssignment: {
			name: "UpdateAssignment",
			I: UpdateAssignmentRequest,
			O: Empty,
			kind: MethodKind.Unary,
		},
		/**
		 * @generated from rpc autograd.v1.AutogradService.DeleteAssignment
		 */
		deleteAssignment: {
			name: "DeleteAssignment",
			I: DeleteByIDRequest,
			O: Empty,
			kind: MethodKind.Unary,
		},
		/**
		 * @generated from rpc autograd.v1.AutogradService.CreateSubmission
		 */
		createSubmission: {
			name: "CreateSubmission",
			I: CreateSubmissionRequest,
			O: CreatedResponse,
			kind: MethodKind.Unary,
		},
		/**
		 * @generated from rpc autograd.v1.AutogradService.UpdateSubmission
		 */
		updateSubmission: {
			name: "UpdateSubmission",
			I: UpdateSubmissionRequest,
			O: Empty,
			kind: MethodKind.Unary,
		},
		/**
		 * @generated from rpc autograd.v1.AutogradService.DeleteSubmission
		 */
		deleteSubmission: {
			name: "DeleteSubmission",
			I: DeleteByIDRequest,
			O: Empty,
			kind: MethodKind.Unary,
		},
		/**
		 * Admin Courses
		 *
		 * @generated from rpc autograd.v1.AutogradService.FindAllAdminCourses
		 */
		findAllAdminCourses: {
			name: "FindAllAdminCourses",
			I: FindAllPaginationRequest,
			O: FindAllAdminCoursesResponse,
			kind: MethodKind.Unary,
		},
		/**
		 * @generated from rpc autograd.v1.AutogradService.CreateAdminCourse
		 */
		createAdminCourse: {
			name: "CreateAdminCourse",
			I: CreateAdminCourseRequest,
			O: CreatedResponse,
			kind: MethodKind.Unary,
		},
		/**
		 * @generated from rpc autograd.v1.AutogradService.UpdateAdminCourse
		 */
		updateAdminCourse: {
			name: "UpdateAdminCourse",
			I: UpdateAdminCourseRequest,
			O: Empty,
			kind: MethodKind.Unary,
		},
		/**
		 * @generated from rpc autograd.v1.AutogradService.FindAdminCourseDetail
		 */
		findAdminCourseDetail: {
			name: "FindAdminCourseDetail",
			I: FindByIDRequest,
			O: FindAdminCourseDetailResponse,
			kind: MethodKind.Unary,
		},
		/**
		 * Admin Course Students
		 *
		 * @generated from rpc autograd.v1.AutogradService.FindAllCourseStudents
		 */
		findAllCourseStudents: {
			name: "FindAllCourseStudents",
			I: FindAllCourseStudentsRequest,
			O: FindAllCourseStudentsResponse,
			kind: MethodKind.Unary,
		},
		/**
		 * Student Assignment
		 * Student Assignment Queries
		 *
		 * @generated from rpc autograd.v1.AutogradService.FindAllStudentAssignments
		 */
		findAllStudentAssignments: {
			name: "FindAllStudentAssignments",
			I: FindAllStudentAssignmentsRequest,
			O: FindAllStudentAssignmentsResponse,
			kind: MethodKind.Unary,
		},
		/**
		 * @generated from rpc autograd.v1.AutogradService.FindStudentAssignment
		 */
		findStudentAssignment: {
			name: "FindStudentAssignment",
			I: FindByIDRequest,
			O: StudentAssignment,
			kind: MethodKind.Unary,
		},
		/**
		 * Student Assignment Command
		 *
		 * @generated from rpc autograd.v1.AutogradService.SubmitStudentSubmission
		 */
		submitStudentSubmission: {
			name: "SubmitStudentSubmission",
			I: SubmitStudentSubmissionRequest,
			O: CreatedResponse,
			kind: MethodKind.Unary,
		},
		/**
		 * @generated from rpc autograd.v1.AutogradService.ResubmitStudentSubmission
		 */
		resubmitStudentSubmission: {
			name: "ResubmitStudentSubmission",
			I: ResubmitStudentSubmissionRequest,
			O: Empty,
			kind: MethodKind.Unary,
		},
		/**
		 * Student Courses
		 *
		 * @generated from rpc autograd.v1.AutogradService.FindAllStudentEnrolledCourses
		 */
		findAllStudentEnrolledCourses: {
			name: "FindAllStudentEnrolledCourses",
			I: FindAllStudentEnrolledCoursesRequest,
			O: FindAllStudentEnrolledCoursesResponse,
			kind: MethodKind.Unary,
		},
		/**
		 * @generated from rpc autograd.v1.AutogradService.FindStudentCourseDetail
		 */
		findStudentCourseDetail: {
			name: "FindStudentCourseDetail",
			I: FindByIDRequest,
			O: FindStudentCourseDetailResponse,
			kind: MethodKind.Unary,
		},
		/**
		 * Auth
		 * Auth Mutation
		 *
		 * @generated from rpc autograd.v1.AutogradService.Login
		 */
		login: {
			name: "Login",
			I: LoginRequest,
			O: LoginResponse,
			kind: MethodKind.Unary,
		},
	},
} as const;
