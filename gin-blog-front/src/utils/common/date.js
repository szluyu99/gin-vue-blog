import dayjs from 'dayjs'

/**
 * 格式化时间，默认格式：YYYY-MM-DD HH:mm:ss
 * @param {String | Date | undefined} time
 * @param {String} format
 * @returns {String}
 */
export function formatDateTime(time, format = 'YYYY-MM-DD HH:mm:ss') {
  return dayjs(time).format(format)
}

/**
 * 格式化日期，默认格式：YYYY-MM-DD
 * @param {String | Date | undefined} date
 * @param {String} format
 * @returns {String}
 */
export function formatDate(date, format = 'YYYY-MM-DD') {
  return formatDateTime(date, format)
}
