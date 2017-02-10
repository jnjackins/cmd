package main

import "testing"
import "regexp"

func TestSubGlobal(t *testing.T) {
	input := "The Quick Brown Fox Jumps Over The Lazy Dog"
	want := "hT.e uQ.ick rB.own oF.x uJ.mps vO.er hT.e aL.zy oD.g"
	re := regexp.MustCompile(`([A-Z])([a-z])`)
	got := sub(input, re, `\2\1.`, true)
	if got != want {
		t.Fatalf("got:\n\t%q\nwant\n\t%q", got, want)
	}
}

var sink string

func BenchmarkSubGlobal(b *testing.B) {
	input := "The Quick Brown Fox Jumps Over The Lazy Dog"
	re := regexp.MustCompile(`([A-Z])([a-z])`)
	for i := 0; i < b.N; i++ {
		sink = sub(input, re, `\2\1.`, true)
	}
}

func BenchmarkSubGlobalLarge(b *testing.B) {
	input := `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam nec purus tortor. Donec purus dolor, commodo ut dapibus vitae, ultrices at erat. Etiam viverra, libero eu dictum auctor, tellus libero pulvinar leo, nec malesuada sapien elit at leo. Morbi porttitor consequat felis, nec sagittis dolor consectetur ut. Ut mollis vel odio posuere scelerisque. Pellentesque tempor, metus sed posuere consequat, eros metus congue sem, et faucibus lorem neque et orci. Cras porta fermentum est vitae convallis. Etiam eleifend massa ac arcu placerat, ut tristique arcu imperdiet. Morbi ut felis dui.

Duis pretium mi eu tellus auctor, ac euismod sem lacinia. Morbi in nibh lobortis, consequat metus et, sodales est. Suspendisse sed vulputate dui. Aliquam in vehicula nibh, viverra vestibulum lectus. Cras ut felis eu nulla pharetra elementum. Pellentesque dignissim nunc ut sem ultricies semper. Cras consectetur porttitor ullamcorper. Nam viverra feugiat felis vitae tempus. Sed iaculis ipsum sit amet ligula rutrum mollis.

Vestibulum aliquam dui eget nulla pretium volutpat. Integer tempor pretium ex, id gravida sem faucibus eu. Aenean vel urna tincidunt, lobortis ex at, eleifend diam. Praesent vehicula, eros nec pulvinar sodales, magna ex ultricies nisi, vel scelerisque elit odio non metus. Sed semper faucibus mi nec vestibulum. Mauris ut orci libero. Duis mi eros, efficitur in orci eget, ultricies mollis arcu. Suspendisse eget sagittis nibh. Morbi imperdiet venenatis interdum. Nullam id nunc malesuada, tempus nisl in, rhoncus nunc. Nullam ultricies ligula nec orci laoreet semper.`
	re := regexp.MustCompile(`([A-Z])([a-z])`)
	for i := 0; i < b.N; i++ {
		sink = sub(input, re, `\2\1.`, true)
	}
}
