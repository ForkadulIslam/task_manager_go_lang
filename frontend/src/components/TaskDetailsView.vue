<template>
  <div v-if="task" class="space-y-6 text-gray-200">
    <!-- Basic Info -->
    <div>
      <h3 class="text-lg font-bold text-white mb-2">{{ task.Label }}</h3>
      <!-- Creator Info -->
      <div v-if="task.Creator && task.Creator.username" class="flex items-center text-gray-400 text-xs mb-4">
        <UserAvatar :username="task.Creator.username" size="sm" class="mr-2" />
        <p>Created by <span class="font-semibold text-gray-300">{{ task.Creator.username }}</span></p>
        <span v-if="task.Creator.user_label" class="ml-2 px-2 py-0.5 rounded-full bg-gray-700 text-gray-300 text-xxs">
          {{ task.Creator.user_label === 1 ? 'Super Admin' : 'User' }}
        </span>
      </div>
      <div class="grid grid-cols-2 gap-4 text-xs">
        <div class="flex items-center">
          <p class="font-semibold text-gray-300 mr-2">Priority:</p>
          <span :class="getPriorityClass(task.Priority)" class="px-2 py-1 text-xs font-semibold rounded-full">{{ task.Priority }}</span>
        </div>
        <div class="flex items-center">
          <p class="font-semibold text-gray-300 flex items-center mr-2"><svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>Status:</p>
          <select v-if="isAssignee" @change="updateTaskStatus" :value="task.Status" class="form-input bg-gray-700 border-gray-600 text-white text-xs py-1 px-2 rounded-full focus:ring-sky-500 focus:border-sky-500">
            <option v-for="status in statusOptions" :key="status" :value="status">{{ status }}</option>
          </select>
          <span v-else :class="getStatusClass(task.Status)" class="px-2 py-1 text-xs font-semibold rounded-full">{{ task.Status }}</span>
        </div>
        <div class="flex items-center">
          <p class="font-semibold text-gray-300 mr-2">Start Date:</p>
          <p>{{ formatDate(task.StartDate) }}</p>
        </div>
        <div class="flex items-center">
          <p class="font-semibold text-gray-300 flex items-center mr-2"><svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>Due Date:</p>
          <p>{{ formatDate(task.DueDate) }}</p>
        </div>
      </div>
    </div>

    <!-- Attachment -->
    <div v-if="task.Attachment">
      <h4 class="font-semibold text-gray-300 mb-2 text-sm">Attachment:</h4>
      <a :href="getAttachmentUrl(task.Attachment)" target="_blank" class="text-sky-400 hover:underline text-xs flex items-center">
        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.5 6H5.25A2.25 2.25 0 003 8.25v10.5A2.25 2.25 0 005.25 21h10.5A2.25 2.25 0 0018 18.75V10.5m-10.5 6L21 3m0 0l-5.25 5.25M21 3H15"></path></svg>
        {{ task.Attachment.split('/').pop().split('-').pop() }}
      </a>
    </div>

    <!-- Assigned Users -->
    <div v-if="task.AssignedUsers && task.AssignedUsers.length">
      <h4 class="font-semibold text-gray-300 mb-2 text-sm">Assigned users:</h4>
      <div class="flex flex-wrap gap-2">
        <span v-for="assignee in task.AssignedUsers" :key="assignee.ID" class="px-3 py-1 rounded-full bg-sky-800 text-sky-200 text-xs font-medium">
          {{ assignee.User.username }}
        </span>
      </div>
    </div>

    <!-- Assigned Groups -->
    <div v-if="task.AssignedGroups && task.AssignedGroups.length">
      <h4 class="font-semibold text-gray-300 mb-2 text-sm">Assigned group:</h4>
      <div class="flex flex-wrap gap-2">
        <span v-for="assignee in task.AssignedGroups" :key="assignee.ID" class="px-3 py-1 rounded-full bg-purple-800 text-purple-200 text-xs font-medium">
          {{ assignee.Group.label }}
          <span v-if="assignee.Group.users && assignee.Group.users.length" class="text-purple-400 ml-1">
            ({{ getUniqueUsernames(assignee.Group.users) }})
          </span>
        </span>
      </div>
    </div>

    <!-- Follow-up Users -->
    <div v-if="task.FollowupUsers && task.FollowupUsers.length">
      <h4 class="font-semibold text-gray-300 mb-2 text-sm">Follow-up Users:</h4>
      <div class="flex flex-wrap gap-2">
        <span v-for="followup in task.FollowupUsers" :key="followup.ID" class="px-3 py-1 rounded-full bg-emerald-800 text-emerald-200 text-xs font-medium">
          {{ followup.User.username }}
        </span>
      </div>
    </div>

    <!-- Follow-up Groups -->
    <div v-if="task.FollowupGroups && task.FollowupGroups.length">
      <h4 class="font-semibold text-gray-300 mb-2 text-sm">Follow-up Groups:</h4>
      <div class="flex flex-wrap gap-2">
        <span v-for="followupGroup in task.FollowupGroups" :key="followupGroup.ID" class="px-3 py-1 rounded-full bg-teal-800 text-teal-200 text-xs font-medium">
          {{ followupGroup.Group.label }}
          <span v-if="followupGroup.Group.users && followupGroup.Group.users.length" class="text-teal-400 ml-1">
            ({{ getUniqueUsernames(followupGroup.Group.users) }})
          </span>
        </span>
      </div>
    </div>

    <!-- Description -->
    <div v-if="task.Description">
      <h4 class="font-semibold text-gray-300 mb-2 text-sm">Description:</h4>
      <div class="prose prose-sm prose-invert max-w-none bg-gray-700 p-3 rounded-md" v-html="task.Description"></div>
    </div>

    <!-- Comments -->
    <div v-if="task.Comments && task.Comments.length">
      <h4 class="font-semibold text-gray-300 mb-2 flex items-center text-sm"><svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z"></path></svg>Comments:</h4>
      <div class="space-y-3 max-h-64 overflow-y-auto pr-2">
        <div v-for="comment in task.Comments" :key="comment.ID" class="bg-gray-700 p-3 rounded-md text-xs">
          <div class="flex items-center mb-2"> <!-- Added flex container -->
            <UserAvatar :username="comment.User.username" size="md" class="mr-2" />
            <p class="font-semibold text-gray-300">{{ comment.User.username }} <span class="text-gray-500 text-xxs ml-2">{{ formatDate(comment.CreatedAt) }}</span></p>
          </div>
          <div class="prose prose-sm prose-invert max-w-none mt-1 text-gray-200" v-html="comment.Comment"></div>
        </div>
    </div>
      </div>
    <!-- Add Comment Section -->
    <div v-if="canComment" class="bg-gray-800 p-4 rounded-md mb-4">
      <h4 class="font-semibold text-gray-300 mb-2 pb-2 border-b border-gray-700 flex items-center text-sm">
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
        Add Comment:
      </h4>
      <div class="mt-4">
        <RichTextEditor v-model="newComment" />
        <button
          @click="submitComment"
          class="mt-2 w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-xs font-medium text-white bg-sky-600 hover:bg-sky-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-sky-500 transition-colors"
        >
          Submit Comment
        </button>
      </div>
    </div>

   
  </div>
  <div v-else class="text-center text-gray-400">
    <p>No task selected or task data is missing.</p>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue';
import apiClient from '../services/api'; // Import apiClient
import UserAvatar from './UserAvatar.vue'; // Import UserAvatar component
import Modal from './Modal.vue'; // Import Modal component
import { useAuthStore } from '../stores/auth'; // Import useAuthStore
import { useToastStore } from '../stores/toast'; // New import
import RichTextEditor from './RichTextEditor.vue'; // Import RichTextEditor

const authStore = useAuthStore();
const authUserID = computed(() => authStore.user?.id);
const toastStore = useToastStore(); // New initialization

const statusOptions = ['Pending', 'In Progress', 'In Review', 'Completed'];

const isAssignee = computed(() => {
  if (!props.task || !authUserID.value) {
    return false;
  }

  const currentUserID = authUserID.value;

  // Check if assigned directly to user
  if (props.task.AssignedUsers && props.task.AssignedUsers.some(au => au.User.id === currentUserID)) {
    return true;
  }

  // Check if assigned via group
  if (props.task.AssignedGroups) {
    for (const ag of props.task.AssignedGroups) {
      if (ag.Group && ag.Group.users && ag.Group.users.some(gu => gu.id === currentUserID)) {
        return true;
      }
    }
  }

  return false;
});

const updateTaskStatus = async (event) => {
  const newStatus = event.target.value;
  try {
    await apiClient.post(`/tasks/${props.task.ID}/status`, { status: newStatus });
    // Update the local task object to reflect the change immediately
    props.task.Status = newStatus;
    toastStore.addToast('Task status updated successfully!', 'success');
  } catch (error) {
    console.error('Failed to update task status:', error);
    toastStore.addToast('Failed to update task status. Please try again.', 'error');
  }
};

const props = defineProps({
  task: {
    type: Object,
    required: true,
  },
});

const emit = defineEmits(['commentSubmitted']);

const newComment = ref(''); // New reactive variable for comment input

const submitComment = async () => { // Make it async
  if (newComment.value.trim() === '') {
    toastStore.addToast('Comment cannot be empty.', 'info');
    return;
  }

  try {
    const taskId = props.task.ID;
    const commentData = {
      Comment: newComment.value,
      // UserID: ... (as discussed, assuming backend handles this from session)
    };

    await apiClient.post(`/tasks/${taskId}/comments`, commentData);

    toastStore.addToast('Comment submitted successfully!', 'success');
    newComment.value = ''; // Clear the textarea

    emit('commentSubmitted'); // Emit a custom event

  } catch (error) {
    console.error('Failed to submit comment:', error);
    toastStore.addToast('Failed to submit comment. Please try again.', 'error');
  }
};

const getUniqueUsernames = (users) => {
  if (!users || !users.length) return '';
  const uniqueUsers = [...new Map(users.map(user => [user.id, user])).values()];
  return uniqueUsers.map(user => user.username).join(', ');
};

const getPriorityClass = (priority) => {
  switch (priority) {
    case 'High': return 'bg-red-500/20 text-red-300';
    case 'Escalation': return 'bg-red-700/40 text-red-200';
    case 'Medium': return 'bg-yellow-500/20 text-yellow-300';
    case 'Normal':
    default:
      return 'bg-sky-500/20 text-sky-300';
  }
};

const getStatusClass = (status) => {
  switch (status) {
    case 'Completed': return 'bg-green-500/20 text-green-300';
    case 'In Progress': return 'bg-blue-500/20 text-blue-300';
    case 'In Review': return 'bg-purple-500/20 text-purple-300';
    case 'Pending':
    default:
      return 'bg-gray-500/20 text-gray-300';
  }
};

const formatDate = (dateString) => {
  if (!dateString) return 'N/A';
  const options = { year: 'numeric', month: 'short', day: 'numeric' };
  return new Date(dateString).toLocaleDateString(undefined, options);
};

const getAttachmentUrl = (path) => {
  // Assuming attachments are served from the base URL + /uploads/
  // You might need to adjust this based on your backend's serving configuration
  return `http://localhost:8080/${path.replace(/\\/g, '/')}`;
};

const canComment = computed(() => {
    //console.log(currentUserID);
  if (!props.task || !authUserID.value) {
    return false;
  }

  const currentUserID = authUserID.value;


  //Check if task created by auth user
  if(props.task.CreatedBy === currentUserID){
    return true;
  }

  // Check if assigned directly to user
  if (props.task.AssignedUsers && props.task.AssignedUsers.some(au => au.UserID === currentUserID)) {
    return true;
  }
  // Check if assigned via group
  if (props.task.AssignedGroups) {
    for (const ag of props.task.AssignedGroups) {
      if (ag.Group && ag.Group.users && ag.Group.users.some(gu => gu.id === currentUserID)) {
        return true;
      }
    }
  }

  console.log(props.task.FollowupUsers);
  // Check if a follow-up user assigned directly
  if (props.task.FollowupUsers && props.task.FollowupUsers.some(fu => fu.UserID === currentUserID)) {
    return true;
  }
  //Check if a follow-up user assigned via group
  if(props.task.FollowupGroups){
    for (const ag of props.task.FollowupGroups) {
      if (ag.Group && ag.Group.users && ag.Group.users.some(gu => gu.id === currentUserID)) {
        return true;
      }
    }
  }

  return false;
});
</script>
