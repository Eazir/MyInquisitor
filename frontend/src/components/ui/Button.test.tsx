import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { Button } from './Button'

describe('Button', () => {
  it('renders children', () => {
    render(<Button>Click me</Button>)
    expect(screen.getByText('Click me')).toBeInTheDocument()
  })

  it('applies variant classes', () => {
    render(<Button variant="danger">Delete</Button>)
    const btn = screen.getByText('Delete')
    expect(btn.className).toContain('danger')
  })

  it('shows loading state', () => {
    render(<Button loading>Loading</Button>)
    expect(screen.getByText('Loading')).toBeDisabled()
  })

  it('calls onClick when clicked', () => {
    let clicked = false
    render(<Button onClick={() => { clicked = true }}>Click</Button>)
    screen.getByText('Click').click()
    expect(clicked).toBe(true)
  })
})
