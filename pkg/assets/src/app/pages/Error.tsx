import { AlertTriangleIcon } from "lucide-react";
import React from "react";

const Error = function () {
  console.log("error");
  return (
    <>
      <div className="w-full p-8">
        <a href="http://www.pointy.ai/" className="scale-[0.85] mt-[-0.25rem]">
          <svg
            className="dark:hidden block"
            width="112"
            height="26"
            viewBox="0 0 112 26"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <g style={{ mixBlendMode: "multiply", opacity: "0.55" }}>
              <path
                d="M5.28907 23.3556L4.71093 12.3708L4.44912 12.3846C4.16855 12.3993 3.89107 12.464 3.6574 12.6199C2.88148 13.1381 1.16055 14.6767 1.3384 18.0559C1.39113 19.0578 1.03251 20.3976 0.647925 21.5143C0.309492 22.4964 1.05169 23.5786 2.08929 23.524L5.28907 23.3556ZM6.6762 17.6621L6.27483 12.1767L6.6762 17.6621ZM13.0234 17.446L12.9163 11.9471L13.0234 17.446ZM26.4032 16.9394L26.4736 22.439L26.4032 16.9394ZM44.4465 17.5242L44.0876 23.0125L44.4465 17.5242ZM76.175 18.3847L76.113 23.8843L76.175 18.3847ZM87.3571 18.5709L87.5761 13.0753L87.3571 18.5709ZM98.1592 18.962L97.9565 24.4583L98.1592 18.962ZM106.98 13.503V24.503H107.242C107.523 24.503 107.803 24.4531 108.045 24.3096C108.847 23.8329 110.646 22.3869 110.646 19.003C110.646 17.9997 111.075 16.6806 111.518 15.5857C111.907 14.6227 111.223 13.503 110.184 13.503H106.98ZM5.28907 23.3556C5.86258 23.3254 6.38442 23.246 6.64438 23.2074C6.80542 23.1835 6.89155 23.1705 6.97678 23.1593C7.0533 23.1491 7.07992 23.1472 7.07756 23.1474L6.27483 12.1767C5.74822 12.2153 5.2462 12.2944 5.02898 12.3266C4.72976 12.3711 4.69337 12.3717 4.71093 12.3708L5.28907 23.3556ZM7.07756 23.1474C9.00268 23.0065 10.9325 22.9878 13.1304 22.945L12.9163 11.9471C10.8831 11.9866 8.57238 12.0086 6.27483 12.1767L7.07756 23.1474ZM13.1304 22.945C15.4589 22.8996 17.7748 22.7826 19.9812 22.675C22.2178 22.566 24.3541 22.4661 26.4736 22.439L26.3327 11.4399C23.9893 11.4699 21.6674 11.5798 19.4458 11.6881C17.1939 11.7978 15.0515 11.9055 12.9163 11.9471L13.1304 22.945ZM26.4736 22.439C32.2957 22.3644 38.1435 22.6239 44.0876 23.0125L44.8053 12.036C38.7268 11.6385 32.5608 11.3601 26.3327 11.4399L26.4736 22.439ZM44.0876 23.0125C54.8427 23.7157 65.5867 23.7657 76.113 23.8843L76.237 12.885C65.5866 12.765 55.1934 12.7152 44.8053 12.036L44.0876 23.0125ZM76.113 23.8843C79.9357 23.9274 83.5127 23.9221 87.138 24.0666L87.5761 13.0753C83.7418 12.9225 79.877 12.926 76.237 12.885L76.113 23.8843ZM87.138 24.0666C90.7411 24.2102 94.4004 24.3271 97.9565 24.4583L98.362 13.4658C94.7174 13.3313 91.1727 13.2186 87.5761 13.0753L87.138 24.0666ZM97.9565 24.4583C100.046 24.5354 102.341 24.6371 104.618 24.6371V13.6371C102.599 13.6371 100.574 13.5474 98.362 13.4658L97.9565 24.4583ZM104.618 24.6371C105.279 24.6371 105.908 24.58 106.257 24.5503C106.688 24.5137 106.854 24.503 106.98 24.503V13.503C106.312 13.503 105.684 13.5594 105.326 13.5898C104.886 13.6272 104.731 13.6371 104.618 13.6371V24.6371Z"
                fill="#8B5CF6"
              ></path>
            </g>
            <path
              d="M7.24798 25V1.15779H17.6023C19.101 1.15779 20.4066 1.41892 21.5192 1.94118C22.6319 2.46344 23.4947 3.20141 24.1078 4.1551C24.7209 5.10879 25.0275 6.23278 25.0275 7.52707V7.93579C25.0275 9.36633 24.6868 10.5244 24.0056 11.4099C23.3244 12.2955 22.4843 12.9427 21.4852 13.3514V13.9645C22.3935 14.0099 23.0974 14.3278 23.5969 14.9182C24.0965 15.4858 24.3462 16.2465 24.3462 17.2002V25H19.8503V17.8473C19.8503 17.3024 19.7027 16.8596 19.4075 16.519C19.135 16.1784 18.6695 16.0081 18.011 16.0081H11.7439V25H7.24798ZM11.7439 11.9208H17.1255C18.1927 11.9208 19.0215 11.637 19.6119 11.0693C20.225 10.479 20.5315 9.70693 20.5315 8.75324V8.41264C20.5315 7.45895 20.2363 6.69827 19.6459 6.1306C19.0555 5.54022 18.2154 5.24503 17.1255 5.24503H11.7439V11.9208ZM34.8993 25.4768C33.219 25.4768 31.7317 25.1249 30.4374 24.421C29.1659 23.6944 28.1668 22.6839 27.4401 21.3896C26.7362 20.0726 26.3843 18.5285 26.3843 16.7574V16.3487C26.3843 14.5776 26.7362 13.0448 27.4401 11.7505C28.1441 10.4335 29.1318 9.42309 30.4034 8.71918C31.675 7.99256 33.1509 7.62925 34.8312 7.62925C36.4888 7.62925 37.9307 8.00391 39.1569 8.75324C40.3831 9.47986 41.3367 10.5017 42.0179 11.8187C42.6992 13.113 43.0398 14.623 43.0398 16.3487V17.8133H30.744C30.7894 18.9713 31.2208 19.9137 32.0383 20.6403C32.8557 21.3669 33.8548 21.7302 35.0356 21.7302C36.2391 21.7302 37.1246 21.4691 37.6923 20.9468C38.26 20.4246 38.6914 19.8455 38.9866 19.2098L42.4948 21.049C42.1769 21.6394 41.7114 22.2865 41.0983 22.9904C40.5079 23.6716 39.7132 24.262 38.7141 24.7616C37.715 25.2384 36.4434 25.4768 34.8993 25.4768ZM30.7781 14.6116H38.68C38.5892 13.6352 38.1918 12.8518 37.4879 12.2615C36.8067 11.6711 35.9098 11.3759 34.7972 11.3759C33.6391 11.3759 32.7195 11.6711 32.0383 12.2615C31.3571 12.8518 30.937 13.6352 30.7781 14.6116ZM48.1999 25L42.8184 8.10609H47.3824L51.2994 21.8665H51.9125L55.8294 8.10609H60.3935L55.0119 25H48.1999ZM61.6146 25V8.10609H65.9062V25H61.6146ZM63.7604 6.1306C62.9884 6.1306 62.3299 5.88082 61.7849 5.38127C61.2626 4.88172 61.0015 4.22322 61.0015 3.40578C61.0015 2.58833 61.2626 1.92983 61.7849 1.43028C62.3299 0.930726 62.9884 0.68095 63.7604 0.68095C64.5551 0.68095 65.2136 0.930726 65.7359 1.43028C66.2581 1.92983 66.5193 2.58833 66.5193 3.40578C66.5193 4.22322 66.2581 4.88172 65.7359 5.38127C65.2136 5.88082 64.5551 6.1306 63.7604 6.1306ZM75.9753 25.4768C73.7727 25.4768 71.9675 25 70.5597 24.0463C69.1519 23.0926 68.3003 21.7302 68.0052 19.9591L71.9561 18.9373C72.1151 19.732 72.3762 20.3564 72.7395 20.8106C73.1256 21.2647 73.591 21.594 74.136 21.7983C74.7037 21.98 75.3168 22.0708 75.9753 22.0708C76.9744 22.0708 77.7123 21.9005 78.1892 21.5599C78.666 21.1966 78.9044 20.7538 78.9044 20.2316C78.9044 19.7093 78.6774 19.3119 78.2232 19.0394C77.7691 18.7443 77.0425 18.5058 76.0434 18.3242L75.0897 18.1539C73.9089 17.9268 72.8304 17.6203 71.854 17.2343C70.8776 16.8255 70.0942 16.2692 69.5038 15.5653C68.9134 14.8614 68.6182 13.9531 68.6182 12.8405C68.6182 11.1602 69.2313 9.87723 70.4575 8.99166C71.6837 8.08339 73.2959 7.62925 75.2941 7.62925C77.1787 7.62925 78.7455 8.04933 79.9944 8.88948C81.2433 9.72964 82.0607 10.8309 82.4467 12.1933L78.4617 13.4195C78.28 12.5566 77.9053 11.9436 77.3377 11.5802C76.7927 11.2169 76.1115 11.0353 75.2941 11.0353C74.4766 11.0353 73.8522 11.1829 73.4207 11.4781C72.9893 11.7505 72.7736 12.1366 72.7736 12.6361C72.7736 13.1811 73.0007 13.5898 73.4548 13.8623C73.9089 14.1121 74.522 14.3051 75.2941 14.4413L76.2477 14.6116C77.5193 14.8387 78.666 15.1452 79.6878 15.5312C80.7324 15.8945 81.5498 16.4282 82.1402 17.1321C82.7533 17.8133 83.0598 18.7443 83.0598 19.925C83.0598 21.6962 82.4127 23.0699 81.1184 24.0463C79.8468 25 78.1324 25.4768 75.9753 25.4768ZM93.1557 25.4768C91.4754 25.4768 89.9654 25.1362 88.6257 24.455C87.286 23.7738 86.2302 22.7861 85.4581 21.4918C84.6861 20.1975 84.3001 18.6421 84.3001 16.8255V16.2806C84.3001 14.464 84.6861 12.9086 85.4581 11.6143C86.2302 10.32 87.286 9.33227 88.6257 8.65106C89.9654 7.96985 91.4754 7.62925 93.1557 7.62925C94.8361 7.62925 96.3461 7.96985 97.6858 8.65106C99.0255 9.33227 100.081 10.32 100.853 11.6143C101.625 12.9086 102.011 14.464 102.011 16.2806V16.8255C102.011 18.6421 101.625 20.1975 100.853 21.4918C100.081 22.7861 99.0255 23.7738 97.6858 24.455C96.3461 25.1362 94.8361 25.4768 93.1557 25.4768ZM93.1557 21.6621C94.4727 21.6621 95.5627 21.242 96.4255 20.4019C97.2884 19.539 97.7198 18.3128 97.7198 16.7233V16.3827C97.7198 14.7933 97.2884 13.5784 96.4255 12.7383C95.5854 11.8754 94.4955 11.444 93.1557 11.444C91.8388 11.444 90.7488 11.8754 89.886 12.7383C89.0231 13.5784 88.5917 14.7933 88.5917 16.3827V16.7233C88.5917 18.3128 89.0231 19.539 89.886 20.4019C90.7488 21.242 91.8388 21.6621 93.1557 21.6621Z"
              fill="#18181B"
            ></path>
          </svg>
          <svg
            className="dark:block hidden"
            width="112"
            height="26"
            viewBox="0 0 112 26"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              opacity="0.55"
              d="M5.28907 23.3556L4.71093 12.3708L4.44912 12.3846C4.16855 12.3993 3.89107 12.464 3.6574 12.6199C2.88148 13.1381 1.16055 14.6767 1.3384 18.0559C1.39113 19.0578 1.03251 20.3976 0.647925 21.5143C0.309492 22.4964 1.05169 23.5786 2.08929 23.524L5.28907 23.3556ZM6.6762 17.6621L6.27483 12.1767L6.6762 17.6621ZM13.0234 17.446L12.9163 11.9471L13.0234 17.446ZM26.4032 16.9394L26.4736 22.439L26.4032 16.9394ZM44.4465 17.5242L44.0876 23.0125L44.4465 17.5242ZM76.175 18.3847L76.113 23.8843L76.175 18.3847ZM87.3571 18.5709L87.5761 13.0753L87.3571 18.5709ZM98.1592 18.962L97.9565 24.4583L98.1592 18.962ZM106.98 13.503V24.503H107.242C107.523 24.503 107.803 24.4531 108.045 24.3096C108.847 23.8329 110.646 22.3869 110.646 19.003C110.646 17.9997 111.075 16.6806 111.518 15.5857C111.907 14.6227 111.223 13.503 110.184 13.503H106.98ZM5.28907 23.3556C5.86258 23.3254 6.38442 23.246 6.64438 23.2074C6.80542 23.1835 6.89155 23.1705 6.97678 23.1593C7.0533 23.1491 7.07992 23.1472 7.07756 23.1474L6.27483 12.1767C5.74822 12.2153 5.2462 12.2944 5.02898 12.3266C4.72976 12.3711 4.69337 12.3717 4.71093 12.3708L5.28907 23.3556ZM7.07756 23.1474C9.00268 23.0065 10.9325 22.9878 13.1304 22.945L12.9163 11.9471C10.8831 11.9866 8.57238 12.0086 6.27483 12.1767L7.07756 23.1474ZM13.1304 22.945C15.4589 22.8996 17.7748 22.7826 19.9812 22.675C22.2178 22.566 24.3541 22.4661 26.4736 22.439L26.3327 11.4399C23.9893 11.4699 21.6674 11.5798 19.4458 11.6881C17.1939 11.7978 15.0515 11.9055 12.9163 11.9471L13.1304 22.945ZM26.4736 22.439C32.2957 22.3644 38.1435 22.6239 44.0876 23.0125L44.8053 12.036C38.7268 11.6385 32.5608 11.3601 26.3327 11.4399L26.4736 22.439ZM44.0876 23.0125C54.8427 23.7157 65.5867 23.7657 76.113 23.8843L76.237 12.885C65.5866 12.765 55.1934 12.7152 44.8053 12.036L44.0876 23.0125ZM76.113 23.8843C79.9357 23.9274 83.5127 23.9221 87.138 24.0666L87.5761 13.0753C83.7418 12.9225 79.877 12.926 76.237 12.885L76.113 23.8843ZM87.138 24.0666C90.7411 24.2102 94.4004 24.3271 97.9565 24.4583L98.362 13.4658C94.7174 13.3313 91.1727 13.2186 87.5761 13.0753L87.138 24.0666ZM97.9565 24.4583C100.046 24.5354 102.341 24.6371 104.618 24.6371V13.6371C102.599 13.6371 100.574 13.5474 98.362 13.4658L97.9565 24.4583ZM104.618 24.6371C105.279 24.6371 105.908 24.58 106.257 24.5503C106.688 24.5137 106.854 24.503 106.98 24.503V13.503C106.312 13.503 105.684 13.5594 105.326 13.5898C104.886 13.6272 104.731 13.6371 104.618 13.6371V24.6371Z"
              fill="#8B5CF6"
            ></path>
            <path
              d="M7.24798 25V1.15779H17.6023C19.101 1.15779 20.4066 1.41892 21.5192 1.94118C22.6319 2.46344 23.4947 3.20141 24.1078 4.1551C24.7209 5.10879 25.0275 6.23278 25.0275 7.52707V7.93579C25.0275 9.36633 24.6868 10.5244 24.0056 11.4099C23.3244 12.2955 22.4843 12.9427 21.4852 13.3514V13.9645C22.3935 14.0099 23.0974 14.3278 23.5969 14.9182C24.0965 15.4858 24.3462 16.2465 24.3462 17.2002V25H19.8503V17.8473C19.8503 17.3024 19.7027 16.8596 19.4075 16.519C19.135 16.1784 18.6695 16.0081 18.011 16.0081H11.7439V25H7.24798ZM11.7439 11.9208H17.1255C18.1927 11.9208 19.0215 11.637 19.6119 11.0693C20.225 10.479 20.5315 9.70693 20.5315 8.75324V8.41264C20.5315 7.45895 20.2363 6.69827 19.6459 6.1306C19.0555 5.54022 18.2154 5.24503 17.1255 5.24503H11.7439V11.9208ZM34.8993 25.4768C33.219 25.4768 31.7317 25.1249 30.4374 24.421C29.1659 23.6944 28.1668 22.6839 27.4401 21.3896C26.7362 20.0726 26.3843 18.5285 26.3843 16.7574V16.3487C26.3843 14.5776 26.7362 13.0448 27.4401 11.7505C28.1441 10.4335 29.1318 9.42309 30.4034 8.71918C31.675 7.99256 33.1509 7.62925 34.8312 7.62925C36.4888 7.62925 37.9307 8.00391 39.1569 8.75324C40.3831 9.47986 41.3367 10.5017 42.0179 11.8187C42.6992 13.113 43.0398 14.623 43.0398 16.3487V17.8133H30.744C30.7894 18.9713 31.2208 19.9137 32.0383 20.6403C32.8557 21.3669 33.8548 21.7302 35.0356 21.7302C36.2391 21.7302 37.1246 21.4691 37.6923 20.9468C38.26 20.4246 38.6914 19.8455 38.9866 19.2098L42.4948 21.049C42.1769 21.6394 41.7114 22.2865 41.0983 22.9904C40.5079 23.6716 39.7132 24.262 38.7141 24.7616C37.715 25.2384 36.4434 25.4768 34.8993 25.4768ZM30.7781 14.6116H38.68C38.5892 13.6352 38.1918 12.8518 37.4879 12.2615C36.8067 11.6711 35.9098 11.3759 34.7972 11.3759C33.6391 11.3759 32.7195 11.6711 32.0383 12.2615C31.3571 12.8518 30.937 13.6352 30.7781 14.6116ZM48.1999 25L42.8184 8.10609H47.3824L51.2994 21.8665H51.9125L55.8294 8.10609H60.3935L55.0119 25H48.1999ZM61.6146 25V8.10609H65.9062V25H61.6146ZM63.7604 6.1306C62.9884 6.1306 62.3299 5.88082 61.7849 5.38127C61.2626 4.88172 61.0015 4.22322 61.0015 3.40578C61.0015 2.58833 61.2626 1.92983 61.7849 1.43028C62.3299 0.930726 62.9884 0.68095 63.7604 0.68095C64.5551 0.68095 65.2136 0.930726 65.7359 1.43028C66.2581 1.92983 66.5193 2.58833 66.5193 3.40578C66.5193 4.22322 66.2581 4.88172 65.7359 5.38127C65.2136 5.88082 64.5551 6.1306 63.7604 6.1306ZM75.9753 25.4768C73.7727 25.4768 71.9675 25 70.5597 24.0463C69.1519 23.0926 68.3003 21.7302 68.0052 19.9591L71.9561 18.9373C72.1151 19.732 72.3762 20.3564 72.7395 20.8106C73.1256 21.2647 73.591 21.594 74.136 21.7983C74.7037 21.98 75.3168 22.0708 75.9753 22.0708C76.9744 22.0708 77.7123 21.9005 78.1892 21.5599C78.666 21.1966 78.9044 20.7538 78.9044 20.2316C78.9044 19.7093 78.6774 19.3119 78.2232 19.0394C77.7691 18.7443 77.0425 18.5058 76.0434 18.3242L75.0897 18.1539C73.9089 17.9268 72.8304 17.6203 71.854 17.2343C70.8776 16.8255 70.0942 16.2692 69.5038 15.5653C68.9134 14.8614 68.6182 13.9531 68.6182 12.8405C68.6182 11.1602 69.2313 9.87723 70.4575 8.99166C71.6837 8.08339 73.2959 7.62925 75.2941 7.62925C77.1787 7.62925 78.7455 8.04933 79.9944 8.88948C81.2433 9.72964 82.0607 10.8309 82.4467 12.1933L78.4617 13.4195C78.28 12.5566 77.9053 11.9436 77.3377 11.5802C76.7927 11.2169 76.1115 11.0353 75.2941 11.0353C74.4766 11.0353 73.8522 11.1829 73.4207 11.4781C72.9893 11.7505 72.7736 12.1366 72.7736 12.6361C72.7736 13.1811 73.0007 13.5898 73.4548 13.8623C73.9089 14.1121 74.522 14.3051 75.2941 14.4413L76.2477 14.6116C77.5193 14.8387 78.666 15.1452 79.6878 15.5312C80.7324 15.8945 81.5498 16.4282 82.1402 17.1321C82.7533 17.8133 83.0598 18.7443 83.0598 19.925C83.0598 21.6962 82.4127 23.0699 81.1184 24.0463C79.8468 25 78.1324 25.4768 75.9753 25.4768ZM93.1557 25.4768C91.4754 25.4768 89.9654 25.1362 88.6257 24.455C87.286 23.7738 86.2302 22.7861 85.4581 21.4918C84.6861 20.1975 84.3001 18.6421 84.3001 16.8255V16.2806C84.3001 14.464 84.6861 12.9086 85.4581 11.6143C86.2302 10.32 87.286 9.33227 88.6257 8.65106C89.9654 7.96985 91.4754 7.62925 93.1557 7.62925C94.8361 7.62925 96.3461 7.96985 97.6858 8.65106C99.0255 9.33227 100.081 10.32 100.853 11.6143C101.625 12.9086 102.011 14.464 102.011 16.2806V16.8255C102.011 18.6421 101.625 20.1975 100.853 21.4918C100.081 22.7861 99.0255 23.7738 97.6858 24.455C96.3461 25.1362 94.8361 25.4768 93.1557 25.4768ZM93.1557 21.6621C94.4727 21.6621 95.5627 21.242 96.4255 20.4019C97.2884 19.539 97.7198 18.3128 97.7198 16.7233V16.3827C97.7198 14.7933 97.2884 13.5784 96.4255 12.7383C95.5854 11.8754 94.4955 11.444 93.1557 11.444C91.8388 11.444 90.7488 11.8754 89.886 12.7383C89.0231 13.5784 88.5917 14.7933 88.5917 16.3827V16.7233C88.5917 18.3128 89.0231 19.539 89.886 20.4019C90.7488 21.242 91.8388 21.6621 93.1557 21.6621Z"
              fill="#FAFAFA"
            ></path>
          </svg>
        </a>
      </div>
      <section className="flex flex-col items-center justify-center h-[calc(100dvh-6.2rem)] text-center">
        <div className="mb-4">
          <div className="w-12 h-12 bg-rose-100 rounded-full flex items-center justify-center">
            <AlertTriangleIcon className="w-6 h-6 stroke-rose-500" />
          </div>
        </div>
        <p className="text-[1rem] leading-1.5rem] font-semibold mb-4">
          This page couldn’t be loaded due to an error.
          <br />
          Our team has been notified and will investigate.
        </p>
        <p className="mt-2 text-center text-foreground text-xs">
          Still having issues? Email us at{" "}
          <a href="mailto:taylor@pointy.ai">taylor@pointy.ai</a>
        </p>
      </section>
    </>
  );
};

export default Error;
