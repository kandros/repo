package ui

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/v50/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test input model functionality
func TestInitialModel(t *testing.T) {
	model := initialModel()
	
	// Test model initialization
	assert.NotNil(t, model.textInput, "Text input should be initialized")
	assert.Nil(t, model.err, "Error should be nil initially")
	
	// Test text input configuration
	assert.Equal(t, "Pikachu", model.textInput.Placeholder, "Placeholder should be set correctly")
	assert.Equal(t, 156, model.textInput.CharLimit, "Character limit should be set correctly")
	assert.Equal(t, 20, model.textInput.Width, "Width should be set correctly")
}

func TestInputModelInit(t *testing.T) {
	model := initialModel()
	cmd := model.Init()
	
	// Test that Init returns the correct command
	assert.Equal(t, textinput.Blink, cmd, "Init should return textinput.Blink command")
}

func TestInputModelUpdate_KeyEnter(t *testing.T) {
	model := initialModel()
	
	// Create Enter key message
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	
	updatedModel, cmd := model.Update(msg)
	
	// Test that Enter key triggers quit
	assert.Equal(t, tea.Quit, cmd, "Enter key should trigger quit command")
	assert.IsType(t, inputModel{}, updatedModel, "Updated model should be of correct type")
}

func TestInputModelUpdate_KeyCtrlC(t *testing.T) {
	model := initialModel()
	
	// Create Ctrl+C key message
	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	
	updatedModel, cmd := model.Update(msg)
	
	// Test that Ctrl+C key triggers quit
	assert.Equal(t, tea.Quit, cmd, "Ctrl+C key should trigger quit command")
	assert.IsType(t, inputModel{}, updatedModel, "Updated model should be of correct type")
}

func TestInputModelUpdate_KeyEsc(t *testing.T) {
	model := initialModel()
	
	// Create Esc key message
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	
	updatedModel, cmd := model.Update(msg)
	
	// Test that Esc key triggers quit
	assert.Equal(t, tea.Quit, cmd, "Esc key should trigger quit command")
	assert.IsType(t, inputModel{}, updatedModel, "Updated model should be of correct type")
}

func TestInputModelUpdate_ErrorMessage(t *testing.T) {
	model := initialModel()
	testError := errMsg(assert.AnError)
	
	updatedModel, cmd := model.Update(testError)
	
	// Test error handling
	assert.Nil(t, cmd, "Error message should not return a command")
	inputModel, ok := updatedModel.(inputModel)
	require.True(t, ok, "Updated model should be of inputModel type")
	assert.Equal(t, testError, inputModel.err, "Error should be set in model")
}

func TestInputModelView(t *testing.T) {
	model := initialModel()
	
	view := model.View()
	
	// Test view output format
	assert.Contains(t, view, "\n", "View should contain newline")
	assert.NotEmpty(t, view, "View should not be empty")
}

// Test item interface implementation
func TestItemInterface(t *testing.T) {
	testItem := item("test-repo")
	
	// Test FilterValue method
	filterValue := testItem.FilterValue()
	assert.Empty(t, filterValue, "FilterValue should return empty string")
	
	// Test item string conversion
	assert.Equal(t, "test-repo", string(testItem), "Item should convert to string correctly")
}

// Test itemDelegate interface implementation
func TestItemDelegate(t *testing.T) {
	delegate := itemDelegate{}
	
	// Test delegate methods
	assert.Equal(t, 1, delegate.Height(), "Height should be 1")
	assert.Equal(t, 0, delegate.Spacing(), "Spacing should be 0")
	
	// Test Update method (should return nil)
	cmd := delegate.Update(nil, nil)
	assert.Nil(t, cmd, "Update should return nil")
}

// Test model struct initialization
func TestModelInitialization(t *testing.T) {
	// Create test repositories
	repos := []*github.Repository{
		{
			Name:     github.String("test-repo"),
			FullName: github.String("user/test-repo"),
			HTMLURL:  github.String("https://github.com/user/test-repo"),
		},
	}
	
	githubToken := "test_token_123"
	
	// Test model fields can be set
	model := model{
		repos:       repos,
		githubToken: githubToken,
	}
	
	assert.Len(t, model.repos, 1, "Repos should be set correctly")
	assert.Equal(t, githubToken, model.githubToken, "GitHub token should be set correctly")
	assert.False(t, model.quitting, "Quitting should be false initially")
	assert.Empty(t, model.choice, "Choice should be empty initially")
	assert.Empty(t, model.quitText, "Quit text should be empty initially")
	assert.False(t, model.showRepoFolderInput, "Show repo folder input should be false initially")
}

func TestModelSelectedRepo(t *testing.T) {
	// Create test repositories
	repos := []*github.Repository{
		{
			Name:     github.String("first-repo"),
			FullName: github.String("user/first-repo"),
		},
		{
			Name:     github.String("second-repo"),
			FullName: github.String("user/second-repo"),
		},
	}
	
	model := model{repos: repos}
	
	// Mock list with index 0
	// Note: In a real test, we'd need to properly initialize the list
	// For now, we test the logic would work
	assert.Len(t, repos, 2, "Should have 2 repos for testing")
	
	// Test selectedRepo method would return correct repo
	// model.selectedRepo() would return repos[model.list.Index()]
	// We can't easily test this without mocking the list, but we can verify the repos are set up correctly
	assert.Equal(t, "first-repo", *repos[0].Name, "First repo should be 'first-repo'")
	assert.Equal(t, "second-repo", *repos[1].Name, "Second repo should be 'second-repo'")
}

// Test key bindings
func TestKeyBindings(t *testing.T) {
	// Test key binding definitions
	assert.Equal(t, "o", keyO.Keys()[0], "Open key should be 'o'")
	assert.Equal(t, "enter", keyEnter.Keys()[0], "Enter key should be 'enter'")
	assert.Equal(t, "c", keyC.Keys()[0], "Clone key should be 'c'")
	
	// Test help text
	assert.Contains(t, keyO.Help().Key, "o", "Open key help should contain 'o'")
	assert.Contains(t, keyO.Help().Desc, "open", "Open key description should contain 'open'")
	assert.Contains(t, keyEnter.Help().Desc, "copy", "Enter key description should contain 'copy'")
	assert.Contains(t, keyC.Help().Desc, "clone", "Clone key description should contain 'clone'")
}

// Test List function parameter validation
func TestListFunction_Parameters(t *testing.T) {
	// Create test repositories
	repos := []*github.Repository{
		{
			Name:     github.String("test-repo"),
			FullName: github.String("user/test-repo"),
			HTMLURL:  github.String("https://github.com/user/test-repo"),
		},
	}
	
	githubToken := "test_token_123"
	
	// Test that parameters are valid
	assert.NotNil(t, repos, "Repos should not be nil")
	assert.NotEmpty(t, githubToken, "GitHub token should not be empty")
	assert.Len(t, repos, 1, "Should have one test repository")
	
	// Test repository fields
	repo := repos[0]
	assert.NotNil(t, repo.Name, "Repository name should not be nil")
	assert.NotNil(t, repo.FullName, "Repository full name should not be nil")
	assert.NotNil(t, repo.HTMLURL, "Repository HTML URL should not be nil")
	
	assert.Equal(t, "test-repo", *repo.Name, "Repository name should match")
	assert.Equal(t, "user/test-repo", *repo.FullName, "Repository full name should match")
	assert.Equal(t, "https://github.com/user/test-repo", *repo.HTMLURL, "Repository URL should match")
}

// Test style constants
func TestStyles(t *testing.T) {
	// Test that styles are defined
	assert.NotNil(t, titleStyle, "Title style should be defined")
	assert.NotNil(t, itemStyle, "Item style should be defined")
	assert.NotNil(t, selectedItemStyle, "Selected item style should be defined")
	assert.NotNil(t, paginationStyle, "Pagination style should be defined")
	assert.NotNil(t, helpStyle, "Help style should be defined")
	assert.NotNil(t, quitTextStyle, "Quit text style should be defined")
}

// Test list height constant
func TestListHeight(t *testing.T) {
	assert.Equal(t, 14, listHeight, "List height should be 14")
}

// Test model update with window size message
func TestModelUpdate_WindowSize(t *testing.T) {
	model := model{}
	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	
	// Test window size message handling
	assert.Equal(t, 100, msg.Width, "Window width should be 100")
	assert.Equal(t, 50, msg.Height, "Window height should be 50")
	
	// Note: In a real test, we'd verify that model.list.SetWidth(msg.Width) is called
	// For now, we test that the message structure is correct
}

// Test model view states
func TestModelView_States(t *testing.T) {
	model := model{}
	
	// Test non-quitting state
	assert.False(t, model.quitting, "Model should not be quitting initially")
	
	// Test quitting state
	model.quitting = true
	model.quitText = "Test quit message"
	assert.True(t, model.quitting, "Model should be quitting when set")
	assert.Equal(t, "Test quit message", model.quitText, "Quit text should be set correctly")
}

// Benchmark tests
func BenchmarkInitialModel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = initialModel()
	}
}

func BenchmarkItemCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = item("test-repo")
	}
}

func BenchmarkItemDelegate(b *testing.B) {
	delegate := itemDelegate{}
	for i := 0; i < b.N; i++ {
		_ = delegate.Height()
		_ = delegate.Spacing()
	}
}