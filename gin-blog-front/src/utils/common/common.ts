import dayjs from 'dayjs'

type Time = undefined | string | Date

/** 格式化时间，默认格式：YYYY-MM-DD HH:mm:ss */
export function formatDateTime(time: Time, format = 'YYYY-MM-DD HH:mm:ss'): string {
  return dayjs(time).format(format)
}

/** 格式化日期，默认格式：YYYY-MM-DD */
export function formatDate(date: Time = undefined, format = 'YYYY-MM-DD') {
  return formatDateTime(date, format)
}
