import { defineStore } from "pinia";

export type CourseData = {
  md5: string,
  sha256: string,
  id: number
};

export type CurrentCustomCourse = {
  courseId?: number,
  data: CourseData[]
};

export const useCurrentCourseStore = defineStore("currentCustomCourse", {
  state: () => ({
    courseId: null,
    data: []
  } as CurrentCustomCourse),

  actions: {
    setter(newValue: CurrentCustomCourse) {
      this.courseId = newValue.courseId;
      this.data = newValue.data;
    }
  }
});
